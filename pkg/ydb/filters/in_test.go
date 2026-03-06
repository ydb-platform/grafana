package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIn(t *testing.T) {
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
			name: "IN keyword in literal (single quotas)",
			in: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
		},
		{
			name: "IN keyword in literal (double quotas)",
			in: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
		},
		{
			name: "IN in lower case",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status in (?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN ?",
				args: []any{1, []any{2, 3}},
			},
		},
		{
			name: "single IN without separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN (?,?,?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN with separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN (?, ?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "IN keyword in multi-line",
			in: args{
				sql: `SELECT * WHERE id IN (
					?,
					?,
					?
				)`,
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  `SELECT * WHERE id IN ?`,
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN with single arg",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN(?)",
				args: []any{1},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1}},
			},
		},
		{
			name: "single IN with and without separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN(?,?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN between args",
			in: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) OR deleted = ?",
				args: []any{"test", 1, 2, 3, "false"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN ? OR deleted = ?",
				args: []any{"test", []any{1, 2, 3}, "false"},
			},
		},
		{
			name: "two IN's between args",
			in: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) AND deleted = ? AND role_id IN (?,?,?) AND role_name LIKE ?",
				args: []any{"test", 1, 2, 3, "false", 4, 5, 6, "admin%"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) AND deleted = ? AND role_id IN (?,?,?) AND role_name LIKE ?",
				args: []any{"test", []any{1, 2, 3}, "false", []any{4, 5, 6}, "admin%"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &IN{}
			sql, args := f.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
