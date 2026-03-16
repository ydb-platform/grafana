package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimit(t *testing.T) {
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
			name: "literal limit",
			in: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT 10",
				args: []any{},
			},
			out: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT 10",
				args: []any{},
			},
		},
		{
			name: "parametrized limit",
			in: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT ?",
				args: []any{10},
			},
			out: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT ?",
				args: []any{uint64(10)},
			},
		},
		{
			name: "parametrized limit with other args",
			in: args{
				sql:  "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT ?",
				args: []any{1, 2, 3, 4},
			},
			out: args{
				sql:  "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT ?",
				args: []any{1, 2, 3, uint64(4)},
			},
		},
		{
			name: "parametrized limit with spaces, tabs and new line",
			in: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT \t\n?",
				args: []any{10},
			},
			out: args{
				sql:  "SELECT * FROM tbl ORDER BY id LIMIT \t\n?",
				args: []any{uint64(10)},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConvertLimitArgToUint64{}
			sql, args := s.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
