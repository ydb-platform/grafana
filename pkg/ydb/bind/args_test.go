package bind

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertPositionalArgsToYdbNamedParameters(t *testing.T) {
	for _, tt := range []struct {
		name    string
		in      string
		args    []any
		exp     string
		expArgs []any
	}{
		{
			name: "happy path",
			in:   "SELECT * WHERE p1=? AND (p2=? OR p3=?)",
			args: []any{1, 2, 3},
			exp:  "SELECT * WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
			expArgs: []any{
				sql.Named("p1", 1),
				sql.Named("p2", 2),
				sql.Named("p3", 3),
			},
		},
		{
			name: "dont convert named arg",
			in:   "SELECT * WHERE p1=? AND (p2=$p2 OR p3=?)",
			args: []any{
				1,
				sql.Named("p2", 2),
				3,
			},
			exp: "SELECT * WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
			expArgs: []any{
				sql.Named("p1", 1),
				sql.Named("p2", 2),
				sql.Named("p3", 3),
			},
		},
		{
			name: "exclude single quotas",
			in:   "SELECT ' IN (?,?,?)' WHERE p1=? AND (p2=? OR p3=?)",
			args: []any{1, 2, 3},
			exp:  "SELECT ' IN (?,?,?)' WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
			expArgs: []any{
				sql.Named("p1", 1),
				sql.Named("p2", 2),
				sql.Named("p3", 3),
			},
		},
		{
			name: "replace nil to NULL",
			in:   "INSERT INTO `tbl` (`id`,`value`,`deleted`) VALUES (?, ?, ?), (?, ?, ?)",
			args: []any{1, 2, nil, 3, 4, true},
			exp:  "INSERT INTO `tbl` (`id`,`value`,`deleted`) VALUES ($p1, $p2, NULL), ($p3, $p4, $p5)",
			expArgs: []any{
				sql.Named("p1", 1),
				sql.Named("p2", 2),
				sql.Named("p3", 3),
				sql.Named("p4", 4),
				sql.Named("p5", true),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertPositionalArgsToYdbNamedParameters{}
			out, outArgs, err := f.Rebind(tt.in, tt.args...)
			require.NoError(t, err)
			require.Equal(t, tt.exp, out)
			require.Len(t, outArgs, len(tt.expArgs))
			for i := range tt.expArgs {
				na, _ := tt.expArgs[i].(sql.NamedArg)
				require.Equal(t, na.Name, outArgs[i].Name)
				require.Equal(t, na.Value, outArgs[i].Value)
			}
		})
	}
}
