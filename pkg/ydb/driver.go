package ydb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/grafana/grafana/pkg/ydb/bind"
	"github.com/grafana/grafana/pkg/ydb/filters"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xslices"
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
	_ driver.Rows               = (*rows)(nil)
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
	rows struct {
		columns    []string
		rows       [][]any
		currentRow int
	}
)

func (r *rows) Columns() []string {
	return r.columns
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	r.currentRow++

	if len(r.rows) >= r.currentRow {
		return io.EOF
	}

	if len(dest) != len(r.rows[r.currentRow]) {
		return xerrors.WithStackTrace(fmt.Errorf("unexpected dest size %d, expected %d", len(dest), len(r.rows[r.currentRow])))
	}

	for i := range dest {
		dest[i] = r.rows[r.currentRow][i]
	}

	return nil
}

var binders = xslices.Transform(filters.Filters, func(f core.Filter) bind.Binder {
	return bind.Func(func(query string, namedArgs ...driver.NamedValue) (string, []driver.NamedValue, error) {
		if ff, has := f.(core.FilterWithArgs); has {
			query, args := ff.DoWithArgs(query, nil, nil, xslices.Transform(namedArgs, func(v driver.NamedValue) any {
				return v.Value
			})...)

			return query, xslices.Transform(args, func(v any) driver.NamedValue {
				if vv, has := v.(driver.NamedValue); has {
					return vv
				}

				if vv, has := v.(sql.NamedArg); has {
					return driver.NamedValue{
						Name:  vv.Name,
						Value: vv.Value,
					}
				}

				return driver.NamedValue{Value: v}
			}), nil
		}

		return f.Do(query, nil, nil), namedArgs, nil
	})
})

func queryContext(ctx context.Context, executor driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "WITH RECURSIVE") {
		return &rows{
			columns: []string{"col"},
			rows: [][]any{
				{1},
				{2},
			},
			currentRow: -1,
		}, nil
	}

	query, args, err := rebind(query, args...)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	r, err := executor.QueryContext(ctx, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func execContext(ctx context.Context, executor driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	query, args, err := rebind(query, args...)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	r, err := executor.ExecContext(ctx, query, args)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return r, nil
}

func rebind(sql string, args ...driver.NamedValue) (_ string, _ []driver.NamedValue, err error) {
	for i, b := range binders {
		sql, args, err = b.Rebind(sql, args...)
		if err != nil {
			return sql, args, xerrors.WithStackTrace(fmt.Errorf("%d rebind failed: %w", i, err))
		}
	}

	return sql, args, nil
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
