package filters

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArgs(t *testing.T) {
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
			name: "happy path",
			in: args{
				sql:  "SELECT * WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql: "SELECT * WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", 2),
					sql.Named("p3", 3),
				},
			},
		},
		{
			name: "exclude single quotas",
			in: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql: "SELECT ' IN (?,?,?)' WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", 2),
					sql.Named("p3", 3),
				},
			},
		},
		{
			name: "exclude double quotas",
			in: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, "2", 3},
			},
			out: args{
				sql: "SELECT \" IN (?,?,?)\" WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", "2"),
					sql.Named("p3", 3),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &Args{}
			sql, args := s.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}

}
