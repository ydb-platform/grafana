package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubstr(t *testing.T) {
	for _, tt := range []struct {
		sql string
		exp string
	}{
		{
			sql: `SELECT SUBSTRING("abcdefg", 3, 1);`,
			exp: `SELECT SUBSTRING("abcdefg", 3, 1);`,
		},
		{
			sql: `SELECT SUBSTR("abcdefg", 3, 1);`,
			exp: `SELECT SUBSTRING("abcdefg", 3, 1);`,
		},
		{
			sql: `SELECT SUBSTR ("abcdefg", 3, 1);`,
			exp: `SELECT SUBSTRING ("abcdefg", 3, 1);`,
		},
		{
			sql: `SELECT "SUBSTR()";`,
			exp: `SELECT "SUBSTR()";`,
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			f := Substr{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
