package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func TestString(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		args []any
		exp  []any
	}{
		{
			name: "string",
			sql:  "SELECT ?, ?, ?",
			args: []any{"a", "b", "c"},
			exp:  []any{"a", "b", "c"},
		},
		{
			name: "datetime64",
			sql:  "SELECT ?, ?, ?",
			args: []any{"2006-01-02 15:04:05", "2007-01-02 15:04:05", "2008-01-02 15:04:05"},
			exp: []any{
				types.Datetime64Value(1136214245),
				types.Datetime64Value(1167750245),
				types.Datetime64Value(1199286245),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertStringToDatetime64{}
			sql, args := f.DoWithArgs(tt.sql, nil, nil, tt.args...)
			require.Equal(t, tt.sql, sql)
			require.Equal(t, tt.exp, args)
		})
	}
}
