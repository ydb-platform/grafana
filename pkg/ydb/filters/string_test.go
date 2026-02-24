package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		args []any
		exp  []any
	}{
		{
			name: "int to int64",
			sql:  "SELECT ?, ?, ?",
			args: []any{"a", "b", "c"},
			exp:  []any{[]byte("a"), []byte("b"), []byte("c")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &String{}
			sql, args := f.DoWithArgs(tt.sql, nil, nil, tt.args...)
			require.Equal(t, tt.sql, sql)
			require.Equal(t, tt.exp, args)
		})
	}
}
