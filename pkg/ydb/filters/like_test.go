package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertLikeToStartsWithEndsWith(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	for _, tt := range []struct {
		name string
		in   args
		out  args
	}{
		{
			name: "StartsWith parametrized LIKE",
			in: args{
				sql:  "SELECT * FROM tbl WHERE role LIKE ?",
				args: []any{"prefix:%"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE StartsWith(role, ?)",
				args: []any{"prefix:"},
			},
		},
		{
			name: "StartsWith parametrized LIKE (column in backticks)",
			in: args{
				sql:  "SELECT * FROM tbl WHERE `role` LIKE ?",
				args: []any{"prefix:%"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE StartsWith(`role`, ?)",
				args: []any{"prefix:"},
			},
		},
		{
			name: "StartsWith literal LIKE",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE 'prefix%'",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE StartsWith(role, 'prefix')",
			},
		},
		{
			name: "StartsWith literal LIKE (double quotas)",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE \"prefix%\"",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE StartsWith(role, \"prefix\")",
			},
		},
		{
			name: "EndsWith parametrized LIKE",
			in: args{
				sql:  "SELECT * FROM tbl WHERE role LIKE ?",
				args: []any{"%suffix"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE EndsWith(role, ?)",
				args: []any{"suffix"},
			},
		},
		{
			name: "EndsWith literal LIKE",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE 'prefix:%'",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE StartsWith(role, 'prefix:')",
			},
		},
		{
			name: "EndsWith literal LIKE (double quotas)",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE \"prefix:%\"",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE StartsWith(role, \"prefix:\")",
			},
		},
		{
			name: "nop parametrized LIKE",
			in: args{
				sql:  "SELECT * FROM tbl WHERE role LIKE ?",
				args: []any{"%middle%"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE role LIKE ?",
				args: []any{"%middle%"},
			},
		},
		{
			name: "nop literal LIKE",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE '%middle%'",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE role LIKE '%middle%'",
			},
		},
		{
			name: "nop literal LIKE (double quotas)",
			in: args{
				sql: "SELECT * FROM tbl WHERE role LIKE \"%middle%\"",
			},
			out: args{
				sql: "SELECT * FROM tbl WHERE role LIKE \"%middle%\"",
			},
		},
		{
			name: "concat like pattern",
			in: args{
				sql:  "SELECT Path FROM `.sys/partition_stats` WHERE Path LIKE '%/' || ?",
				args: []any{"migration_log"},
			},
			out: args{
				sql:  "SELECT Path FROM `.sys/partition_stats` WHERE EndsWith(Path, '/' || ?)",
				args: []any{"migration_log"},
			},
		},
		{
			name: "multi concat like pattern",
			in: args{
				sql:  "SELECT Path FROM `.sys/partition_stats` WHERE Path LIKE '%/' || ? || '/' || ? || '/indexImplTable'",
				args: []any{"migration_log", "UQE_some_index"},
			},
			out: args{
				sql:  "SELECT Path FROM `.sys/partition_stats` WHERE EndsWith(Path, '/' || ? || '/' || ? || '/indexImplTable')",
				args: []any{"migration_log", "UQE_some_index"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertLikeToStartsWithEndsWith{}
			sql, args := f.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
