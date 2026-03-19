package bind

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
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
				driver.NamedValue{Name: "p1", Value: 1},
				driver.NamedValue{Name: "p2", Value: 2},
				driver.NamedValue{Name: "p3", Value: 3},
			},
		},
		{
			name: "dont convert named arg",
			in:   "SELECT * WHERE p1=$p1 OR p2=? AND (p3=$p3 OR p4=?)",
			args: []any{
				sql.Named("p1", 1),
				2,
				sql.Named("p3", 3),
				4,
			},
			exp: "SELECT * WHERE p1=$p1 OR p2=$p2 AND (p3=$p3 OR p4=$p4)",
			expArgs: []any{
				driver.NamedValue{Name: "p1", Value: 1},
				driver.NamedValue{Name: "p2", Value: 2},
				driver.NamedValue{Name: "p3", Value: 3},
				driver.NamedValue{Name: "p4", Value: 4},
			},
		},
		{
			name: "exclude single quotas",
			in:   "SELECT ' IN (?,?,?)' WHERE p1=? AND (p2=? OR p3=?)",
			args: []any{1, 2, 3},
			exp:  "SELECT ' IN (?,?,?)' WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
			expArgs: []any{
				driver.NamedValue{Name: "p1", Value: 1},
				driver.NamedValue{Name: "p2", Value: 2},
				driver.NamedValue{Name: "p3", Value: 3},
			},
		},
		{
			name: "replace nil to NULL",
			in:   "INSERT INTO `tbl` (`id`,`value`,`deleted`) VALUES (?, ?, ?), (?, ?, ?)",
			args: []any{1, 2, nil, 3, 4, true},
			exp:  "INSERT INTO `tbl` (`id`,`value`,`deleted`) VALUES ($p1, $p2, NULL), ($p3, $p4, $p5)",
			expArgs: []any{
				driver.NamedValue{Name: "p1", Value: 1},
				driver.NamedValue{Name: "p2", Value: 2},
				driver.NamedValue{Name: "p3", Value: 3},
				driver.NamedValue{Name: "p4", Value: 4},
				driver.NamedValue{Name: "p5", Value: true},
			},
		},
		{
			name: "dont convert named arg + NULL",
			in:   "SELECT * WHERE p1=$p1 OR p2=? AND (p3=$p3 OR p4=?)",
			args: []any{
				sql.Named("p1", 1),
				nil,
				sql.Named("p3", 3),
				4,
			},
			exp: "SELECT * WHERE p1=$p1 OR p2=NULL AND (p3=$p3 OR p4=$p2)",
			expArgs: []any{
				driver.NamedValue{Name: "p1", Value: 1},
				driver.NamedValue{Name: "p3", Value: 3},
				driver.NamedValue{Name: "p2", Value: 4},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertPositionalArgsToYdbNamedParameters{}
			sql, args, err := f.Rebind(tt.in, tt.args...)
			require.NoError(t, err)
			assert.Equal(t, tt.exp, sql)
			assert.Equal(t, args, tt.expArgs)
		})
	}
}
