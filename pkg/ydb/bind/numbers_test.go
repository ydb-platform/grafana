package bind

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumbers(t *testing.T) {
	makeArgs := func(vals ...any) []driver.NamedValue {
		out := make([]driver.NamedValue, len(vals))
		for i, v := range vals {
			out[i] = driver.NamedValue{Value: v}
		}
		return out
	}
	for _, tt := range []struct {
		name   string
		sql    string
		args   []any
		expArgs []any
	}{
		{
			name:   "int to int64",
			sql:    "SELECT ?, ?, ?",
			args:   []any{1, 2, 3},
			expArgs: []any{int64(1), int64(2), int64(3)},
		},
		{
			name:   "float32 to float64",
			sql:    "SELECT ?, ?, ?",
			args:   []any{float32(1), float32(2), float32(3)},
			expArgs: []any{float64(1), float64(2), float64(3)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertNumbersToInt64{}
			outSQL, outArgs, err := f.Rebind(tt.sql, makeArgs(tt.args...)...)
			require.NoError(t, err)
			require.Equal(t, tt.sql, outSQL)
			require.Len(t, outArgs, len(tt.expArgs))
			for i := range tt.expArgs {
				require.Equal(t, tt.expArgs[i], outArgs[i].Value)
			}
		})
	}
}
