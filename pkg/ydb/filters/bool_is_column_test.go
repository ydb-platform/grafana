package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolIsColumn(t *testing.T) {
	for _, tt := range []struct {
		sql string
		exp string
	}{
		{
			sql: `SELECT * FROM tbl WHERE is_user = 1;`,
			exp: `SELECT * FROM tbl WHERE is_user;`,
		},
		{
			sql: `SELECT * FROM tbl WHERE is_active = 1;`,
			exp: `SELECT * FROM tbl WHERE is_active;`,
		},
		{
			sql: `SELECT * FROM tbl WHERE is_folder = 1;`,
			exp: `SELECT * FROM tbl WHERE is_folder;`,
		},
		{
			sql: `SELECT * FROM tbl WHERE is_folder=1;`,
			exp: `SELECT * FROM tbl WHERE is_folder;`,
		},
		{
			sql: `SELECT * FROM tbl WHERE is_folder= 1;`,
			exp: `SELECT * FROM tbl WHERE is_folder;`,
		},
		{
			sql: `SELECT * FROM tbl WHERE is_folder =1;`,
			exp: `SELECT * FROM tbl WHERE is_folder;`,
		},
		{
			sql: `SELECT * FROM tbl AS t WHERE t.is_user = 1;`,
			exp: `SELECT * FROM tbl AS t WHERE t.is_user;`,
		},
		{
			sql: `SELECT * FROM tbl AS t WHERE t.is_active = 1;`,
			exp: `SELECT * FROM tbl AS t WHERE t.is_active;`,
		},
		{
			sql: `SELECT * FROM tbl AS t WHERE t.is_folder = 1;`,
			exp: `SELECT * FROM tbl AS t WHERE t.is_folder;`,
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			f := BoolIsColumn{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
