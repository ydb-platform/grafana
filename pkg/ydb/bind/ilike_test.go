package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestILike(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "simple column name",
			in:   "SELECT * FROM tbl WHERE name ILIKE ?",
			out:  "SELECT * FROM tbl WHERE Unicode::ToLower(name) LIKE Unicode::ToLower(?)",
		},
		{
			name: "column name with table alias",
			in:   "SELECT * FROM tbl AS t WHERE t.name ILIKE ?",
			out:  "SELECT * FROM tbl AS t WHERE Unicode::ToLower(t.name) LIKE Unicode::ToLower(?)",
		},
		{
			name: "simple column name in backticks",
			in:   "SELECT * FROM tbl WHERE `name` ILIKE ?",
			out:  "SELECT * FROM tbl WHERE Unicode::ToLower(`name`) LIKE Unicode::ToLower(?)",
		},
		{
			name: "no ilike",
			in:   "SELECT * FROM t WHERE id = $1",
			out:  "SELECT * FROM t WHERE id = $1",
		},
		{
			name: "dashboard title ilike",
			in:   "WHERE dashboard.org_id=$1 AND dashboard.title ILIKE $p2",
			out:  "WHERE dashboard.org_id=$1 AND Unicode::ToLower(dashboard.title) LIKE Unicode::ToLower($p2)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertILikeToLikeLowerCase{}
			out, _, err := f.Rebind(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.out, out)
		})
	}
}
