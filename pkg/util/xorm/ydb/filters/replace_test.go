package filters

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
			exp: `SELECT Unicode::ReplaceAll("absence", "enc", "");`,
		},
		{
			sql: `SELECT REPLACE ("absence", "enc", "");`,
			exp: `SELECT Unicode::ReplaceAll ("absence", "enc", "");`,
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
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
