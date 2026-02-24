package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumbers(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		args []any
		exp  []any
	}{
		{
			name: "int to int64",
			sql:  "SELECT ?, ?, ?",
			args: []any{1, 2, 3},
			exp:  []any{int64(1), int64(2), int64(3)},
		},
		{
			name: "float32 to float64",
			sql:  "SELECT ?, ?, ?",
			args: []any{float32(1), float32(2), float32(3)},
			exp:  []any{float64(1), float64(2), float64(3)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &Numbers{}
			sql, args := f.DoWithArgs(tt.sql, nil, nil, tt.args...)
			require.Equal(t, tt.sql, sql)
			require.Equal(t, tt.exp, args)
		})
	}
}
