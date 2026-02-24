package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcatOrgID(t *testing.T) {
	for _, tt := range []struct {
		sql string
		exp string
	}{
		{
			sql: `SELECT 'sa-' || CAST(org_id AS TEXT) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS TEXT) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS TEXT) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS Text) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS Text) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS Text) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS UTF8) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS UTF8) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS UTF8) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS UTF8) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS UTF8) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS UTF8) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS STRING) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS STRING) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS STRING) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS String) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS String) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS String) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS BYTES) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS BYTES) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS BYTES) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS BYTES) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS BYTES) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS BYTES) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS Bytes) FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) FROM tbl;`,
		},
		{
			sql: `SELECT CAST(org_id AS Bytes) || '-org' FROM tbl;`,
			exp: `SELECT Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-' || CAST(org_id AS Bytes) || '-org' FROM tbl;`,
			exp: `SELECT 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-org' FROM tbl;`,
		},
		{
			sql: `SELECT 'sa-'||CAST(org_id AS Bytes)||'sa-' FROM tbl;`,
			exp: `SELECT 'sa-'||Unwrap(CAST(org_id AS STRING))||'sa-' FROM tbl;`,
		},
		{
			sql: `SELECT "'sa-'||CAST(org_id AS Bytes)||'sa-'" FROM tbl;`,
			exp: `SELECT "'sa-'||CAST(org_id AS Bytes)||'sa-'" FROM tbl;`,
		},
		{
			sql: "SELECT `'sa-'||CAST(org_id AS Bytes)||'sa-'` FROM tbl;",
			exp: "SELECT `'sa-'||CAST(org_id AS Bytes)||'sa-'` FROM tbl;",
		},
		{
			sql: `SELECT "sa-" || CAST(org_id AS Bytes) || "-org" FROM tbl;`,
			exp: `SELECT "sa-" || Unwrap(CAST(org_id AS STRING)) || "-org" FROM tbl;`,
		},
		{
			sql: "SELECT `sa-` || CAST(org_id AS Bytes) || `-org` FROM tbl;",
			exp: "SELECT `sa-` || Unwrap(CAST(org_id AS STRING)) || `-org` FROM tbl;",
		},
		{
			sql: "SELECT `sa-` || CAST(t.org_id AS Bytes) || `-org` FROM tbl AS t;",
			exp: "SELECT `sa-` || Unwrap(CAST(t.org_id AS STRING)) || `-org` FROM tbl AS t;",
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			f := ConcatOrgID{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
