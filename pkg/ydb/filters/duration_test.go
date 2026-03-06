package filters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
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
			in: args{
				sql: "SELECT * FROM tbl ORDER BY created_at > ?",
				args: []any{
					time.Second,
				},
			},
			out: args{
				sql: "SELECT * FROM tbl ORDER BY created_at > ?",
				args: []any{
					int64(time.Second),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &Duration{}
			sql, args := f.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
