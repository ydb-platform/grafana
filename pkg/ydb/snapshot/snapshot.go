package snapshot

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xstring"
)

//go:embed schema/*.sql
var snapshotFS embed.FS

type sqlFile struct {
	name    string
	content string
}

//go:embed migration_log_schema.template
var createMigrationLogTableTemplate string

//go:embed migration_log_data.template
var upsertMigrationLogData string

func readSnapshotSQL() ([]sqlFile, error) {
	entries, err := snapshotFS.ReadDir("schema")
	if err != nil {
		return nil, xerrors.WithStackTrace(fmt.Errorf("snapshot: read dir: %v", err))
	}

	sqlFiles := make([]sqlFile, 0)
	for _, e := range entries {
		bb, err := snapshotFS.ReadFile("schema/" + e.Name())
		if err != nil {
			return nil, xerrors.WithStackTrace(fmt.Errorf("snapshot: read %s: %v", e.Name(), err))
		}
		sqlFiles = append(sqlFiles, sqlFile{
			name:    e.Name(),
			content: xstring.FromBytes(bb),
		})
	}

	return sqlFiles, nil
}

func Apply(ctx context.Context, db *ydb.Driver, migrationLogTableName string) error {
	sqlFiles, err := readSnapshotSQL()
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	for _, sqlFile := range sqlFiles {
		if err := db.Query().Exec(ctx, sqlFile.content); err != nil {
			return xerrors.WithStackTrace(fmt.Errorf("snapshot: apply %q failed: %v", sqlFile.name, err))
		}
	}

	tmpl, err := template.New("").Parse(createMigrationLogTableTemplate)
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	var buffer strings.Builder

	if err := tmpl.Execute(&buffer, struct{ Name string }{migrationLogTableName}); err != nil {
		return xerrors.WithStackTrace(err)
	}

	if err := db.Query().Exec(ctx, buffer.String()); err != nil {
		return xerrors.WithStackTrace(fmt.Errorf("snapshot: create table %q failed: %v", migrationLogTableName, err))
	}

	tmpl, err = template.New("").Parse(upsertMigrationLogData)
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	buffer.Reset()

	if err := tmpl.Execute(&buffer, struct{ Name string }{migrationLogTableName}); err != nil {
		return xerrors.WithStackTrace(err)
	}

	if err := db.Query().Exec(ctx, buffer.String()); err != nil {
		return xerrors.WithStackTrace(fmt.Errorf("snapshot: upsert data into table %q failed: %v", migrationLogTableName, err))
	}

	return nil
}
