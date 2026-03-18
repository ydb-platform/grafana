package bind

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertLikeToStartsWithEndsWith(t *testing.T) {
	makeArgs := func(vals ...any) []driver.NamedValue {
		if len(vals) == 0 {
			return nil
		}
		out := make([]driver.NamedValue, len(vals))
		for i, v := range vals {
			out[i] = driver.NamedValue{Value: v}
		}
		return out
	}
	for _, tt := range []struct {
		name   string
		inSQL  string
		inArgs []any
		expSQL string
		expArgs []any
	}{
		{
			name:   "StartsWith parametrized LIKE",
			inSQL:  "SELECT * FROM tbl WHERE role LIKE ?",
			inArgs: []any{"prefix:%"},
			expSQL: "SELECT * FROM tbl WHERE StartsWith(role, ?)",
			expArgs: []any{"prefix:"},
		},
		{
			name:   "EndsWith parametrized LIKE",
			inSQL:  "SELECT * FROM tbl WHERE role LIKE ?",
			inArgs: []any{"%suffix"},
			expSQL: "SELECT * FROM tbl WHERE EndsWith(role, ?)",
			expArgs: []any{"suffix"},
		},
		{
			name:   "nop parametrized LIKE",
			inSQL:  "SELECT * FROM tbl WHERE role LIKE ?",
			inArgs: []any{"%middle%"},
			expSQL: "SELECT * FROM tbl WHERE role LIKE ?",
			expArgs: []any{"%middle%"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertLikeToStartsWithEndsWith{}
			outSQL, outArgs, err := f.Rebind(tt.inSQL, makeArgs(tt.inArgs...)...)
			require.NoError(t, err)
			require.Equal(t, tt.expSQL, outSQL)
			require.Len(t, outArgs, len(tt.expArgs))
			for i := range tt.expArgs {
				require.Equal(t, tt.expArgs[i], outArgs[i].Value)
			}
		})
	}
}
