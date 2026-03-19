package ydb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/grafana/grafana/pkg/ydb/bind"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
)

var (
	_ driver.Connector          = (*connector)(nil)
	_ driver.Conn               = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.Tx                 = (*tx)(nil)
	_ driver.ExecerContext      = (*tx)(nil)
	_ driver.QueryerContext     = (*tx)(nil)
	_ driver.Stmt               = (*stmt)(nil)
	_ driver.StmtQueryContext   = (*stmt)(nil)
	_ driver.StmtExecContext    = (*stmt)(nil)
)

type (
	connector struct {
		c driver.Connector
	}
	extendedConn interface {
		driver.Conn
		driver.ConnPrepareContext
		driver.ConnBeginTx
		driver.ExecerContext
		driver.QueryerContext
	}
	conn struct {
		cc extendedConn
	}
	executor interface {
		driver.ExecerContext
		driver.QueryerContext
	}
	extendedTx interface {
		driver.Tx
		executor
	}
	tx struct {
		tx extendedTx
	}
	stmt struct {
		query    string
		executor executor
	}
)

// queryBinders применяются к запросам, возвращающим строки (DQL: SELECT и т.п.).
// Включают PRAGMA AnsiImplicitCrossJoin и преобразования, специфичные для выборки.
var queryBinders = []bind.Binder{
	bind.NewPatches(bind.PATCH_SELECT),
	&bind.Replace{},
	&bind.Lower{},
	&bind.ReduceDuplicateIdInSelect{},
	&bind.ConvertILikeToLikeLowerCase{},
	&bind.ConvertLikeToStartsWithEndsWith{},
	&bind.ConvertSubstrToSubstring{},
	&bind.ConvertNumbersToInt64{},
	&bind.ConvertStringToDatetime64{},
	&bind.ConvertTimeToDatetime64{},
	&bind.ConvertDurationToInt64{},
	&bind.CastLimitOffsetToUint64{},
	&bind.ConvertInArgsToList{},
	&bind.ConvertPositionalArgsToYdbNamedParameters{},
	bind.Prepend("PRAGMA ydb.CostBasedOptimization = \"on\""),
	bind.Prepend("PRAGMA AnsiImplicitCrossJoin"),
}

// execBinders применяются к запросам без возврата строк (DDL/DML: INSERT, UPDATE, DELETE, ALTER и т.п.).
// Включают DDL-биндинги (AddColumn, CreateIndex, …) и преобразования аргументов, без PRAGMA для DQL.
var execBinders = []bind.Binder{
	bind.NewPatches(bind.PATCH_DELETE),
	bind.NewPatches(bind.PATCH_UPDATE),
	bind.NewPatches(bind.PATCH_INSERT),
	&bind.RemoveOrderByFromUpdate{},
	&bind.Replace{},
	&bind.Lower{},
	&bind.ConvertILikeToLikeLowerCase{},
	&bind.ConvertLikeToStartsWithEndsWith{},
	&bind.ConvertSubstrToSubstring{},
	&bind.ConvertNumbersToInt64{},
	&bind.ConvertStringToDatetime64{},
	&bind.ConvertTimeToDatetime64{},
	&bind.ConvertDurationToInt64{},
	&bind.CastLimitOffsetToUint64{},
	&bind.ConvertInArgsToList{},
	&bind.ConvertPositionalArgsToYdbNamedParameters{},
}

func namedToAny(in []driver.NamedValue) []any {
	out := make([]any, 0, len(in))
	for i := range in {
		if in[i].Name == "" {
			out = append(out, in[i].Value)
		} else {
			out = append(out, sql.Named(in[i].Name, in[i].Value))
		}
	}
	return out
}

func anyToNamed(in []any) []driver.NamedValue {
	out := make([]driver.NamedValue, 0, len(in))
	for i := range in {
		switch t := in[i].(type) {
		case sql.NamedArg:
			out = append(out, driver.NamedValue{
				Name:  t.Name,
				Value: t.Value,
			})
		case driver.NamedValue:
			out = append(out, t)
		default:
			panic(fmt.Sprintf("unexpected arg %d type %T", i, t))
		}
	}
	return out
}

func queryContext(ctx context.Context, executor driver.QueryerContext, query string, args []driver.NamedValue) (_ driver.Rows, err error) {
	if strings.Contains(query, "WITH RECURSIVE") {
		return nil, xerrors.WithStackTrace(&mysql.MySQLError{
			Number: mysqlerr.ER_NOT_SUPPORTED_YET,
		})
	}

	params := namedToAny(args)

	for i, b := range queryBinders {
		query, params, err = b.Rebind(query, params...)
		if err != nil {
			return nil, xerrors.WithStackTrace(fmt.Errorf("%d rebind failed: %w", i, err))
		}
	}

	r, err := executor.QueryContext(ctx, query, anyToNamed(params))
	if err != nil {
		var b strings.Builder
		for _, arg := range args {
			fmt.Fprintf(&b, "$%s = %v;\n", arg.Name, arg.Value)
		}
		fmt.Println(err)
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func execContext(ctx context.Context, executor driver.ExecerContext, query string, args []driver.NamedValue) (_ driver.Result, err error) {
	params := namedToAny(args)

	for i, b := range execBinders {
		query, params, err = b.Rebind(query, params...)
		if err != nil {
			return nil, xerrors.WithStackTrace(fmt.Errorf("%d rebind failed: %w", i, err))
		}
	}

	r, err := executor.ExecContext(ctx, query, anyToNamed(params))
	if err != nil {
		var b strings.Builder
		for _, arg := range args {
			fmt.Fprintf(&b, "$%s = %v;\n", arg.Name, arg.Value)
		}
		fmt.Println(err)
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	r, err := execContext(ctx, c.cc, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	r, err := queryContext(ctx, c.cc, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (t *tx) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	r, err := queryContext(ctx, t.tx, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (t *tx) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	r, err := execContext(ctx, t.tx, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (stmt *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	r, err := execContext(ctx, stmt.executor, stmt.query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (stmt *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	r, err := queryContext(ctx, stmt.executor, stmt.query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func (stmt *stmt) Close() error {
	return nil
}

func (stmt *stmt) NumInput() int {
	return -1
}

func (stmt *stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, xerrors.WithStackTrace(driver.ErrSkip)
}

func (stmt *stmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, xerrors.WithStackTrace(driver.ErrSkip)
}

func (t *tx) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (t *tx) Rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (c *connector) Open(dsn string) (driver.Conn, error) {
	// ignore dsn, using parent connector with their options to connect
	cc, err := c.Connect(context.Background())
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	if ccc, has := cc.(extendedConn); has {
		return &conn{cc: ccc}, nil
	}

	_ = cc.Close()

	return nil, xerrors.WithStackTrace(fmt.Errorf("conn is not extended conn: %T", cc))
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	cc, err := c.c.Connect(ctx)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	if ccc, has := cc.(extendedConn); has {
		return &conn{cc: ccc}, nil
	}

	_ = cc.Close()

	return nil, xerrors.WithStackTrace(fmt.Errorf("conn is not extended conn: %T", cc))
}

func (c *connector) Driver() driver.Driver {
	return c
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	switch opts.Isolation {
	case driver.IsolationLevel(sql.LevelSnapshot), driver.IsolationLevel(sql.LevelSerializable):
		// nop
	default:
		if opts.ReadOnly {
			opts.Isolation = driver.IsolationLevel(sql.LevelSnapshot)
		} else {
			opts.Isolation = driver.IsolationLevel(sql.LevelSerializable)
		}
	}

	t, err := c.cc.BeginTx(ctx, opts)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	if tt, has := t.(extendedTx); has {
		return &tx{
			tx: tt,
		}, nil
	}

	return nil, xerrors.WithStackTrace(fmt.Errorf("tx is not extended tx: %T", t))

}

func (c *conn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	return &stmt{
		query:    query,
		executor: c,
	}, nil
}

func (c *conn) Prepare(string) (driver.Stmt, error) {
	return nil, xerrors.WithStackTrace(driver.ErrSkip)
}

func (c *conn) Close() error {
	err := c.cc.Close()
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	return nil, xerrors.WithStackTrace(driver.ErrSkip)
}
