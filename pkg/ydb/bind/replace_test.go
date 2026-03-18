package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplace(t *testing.T) {
	for _, tt := range []struct {
		sql string
		exp string
	}{
		{
			sql: `SELECT REPLACE("absence", "enc", "");`,
			exp: `SELECT Unicode::ReplaceAll(CAST("absence" AS Text), "enc", "");`,
		},
		{
			sql: `SELECT REPLACE ("absence", "enc", "");`,
			exp: `SELECT Unicode::ReplaceAll (CAST("absence" AS Text), "enc", "");`,
		},
		{
			sql: `SELECT REPLACE ('absence', 'enc', '');`,
			exp: `SELECT Unicode::ReplaceAll (CAST('absence' AS Text), 'enc', '');`,
		},
		{
			sql: `SELECT REPLACE ('lorem, ipsum', 'enc', '');`,
			exp: `SELECT Unicode::ReplaceAll (CAST('lorem, ipsum' AS Text), 'enc', '');`,
		},
		{
			sql: `SELECT REPLACE ("lorem, ipsum", 'enc', '');`,
			exp: `SELECT Unicode::ReplaceAll (CAST("lorem, ipsum" AS Text), 'enc', '');`,
		},
		{
			sql: "SELECT REPLACE (`lorem, ipsum`, 'enc', '');",
			exp: "SELECT Unicode::ReplaceAll (CAST(`lorem, ipsum` AS Text), 'enc', '');",
		},
		{
			sql: `SELECT REPLACE ('absence'u, 'enc'u, ''u);`,
			exp: `SELECT Unicode::ReplaceAll (CAST('absence'u AS Text), 'enc'u, ''u);`,
		},
		{
			sql: "SELECT REPLACE(`absence`, `enc`, ``) FROM tbl;",
			exp: "SELECT Unicode::ReplaceAll(CAST(`absence` AS Text), `enc`, ``) FROM tbl;",
		},
		{
			sql: `SELECT REPLACE(col1, "enc", "") FROM tbl;`,
			exp: `SELECT Unicode::ReplaceAll(CAST(col1 AS Text), "enc", "") FROM tbl;`,
		},
		{
			sql: `SELECT REPLACE(t.col1, "enc", "") FROM tbl AS t;`,
			exp: `SELECT Unicode::ReplaceAll(CAST(t.col1 AS Text), "enc", "") FROM tbl AS t;`,
		},
		{
			sql: "SELECT REPLACE(`t`.`col1`, \"enc\", \"\") FROM `tbl` AS `t`;",
			exp: "SELECT Unicode::ReplaceAll(CAST(`t`.`col1` AS Text), \"enc\", \"\") FROM `tbl` AS `t`;",
		},
		{
			sql: `SELECT REPLACED_COLUMN FROM tbl;`,
			exp: `SELECT REPLACED_COLUMN FROM tbl;`,
		},
		{
			sql: `SELECT COLUMN_REPLACED FROM tbl;`,
			exp: `SELECT COLUMN_REPLACED FROM tbl;`,
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			f := Replace{}
			out, _, err := f.Rebind(tt.sql)
			require.NoError(t, err)
			require.Equal(t, tt.exp, out)
		})
	}
}
