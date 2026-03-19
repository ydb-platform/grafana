package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimit(t *testing.T) {
	for _, tt := range []struct {
		name    string
		inSQL   string
		inArgs  []any
		expSQL  string
		expArgs []any
	}{
		{
			name:    "parametrized limit",
			inSQL:   "SELECT * FROM tbl ORDER BY id LIMIT ?",
			inArgs:  []any{10},
			expSQL:  "SELECT * FROM tbl ORDER BY id LIMIT CAST(? AS Uint64)",
			expArgs: []any{10},
		},
		{
			name:    "parametrized limit+offset",
			inSQL:   "SELECT * FROM tbl ORDER BY id LIMIT ? OFFSET ?",
			inArgs:  []any{10},
			expSQL:  "SELECT * FROM tbl ORDER BY id LIMIT CAST(? AS Uint64) OFFSET CAST(? AS Uint64)",
			expArgs: []any{10},
		},
		{
			name: "multiline limit+offset",
			inSQL: `SELECT *
					FROM tbl
					ORDER BY id
					LIMIT ?
					OFFSET ?`,
			inArgs: []any{10},
			expSQL: `SELECT *
					FROM tbl
					ORDER BY id
					LIMIT CAST(? AS Uint64)
					OFFSET CAST(? AS Uint64)`,
			expArgs: []any{10},
		},
		{
			name:    "parametrized limit with other args",
			inSQL:   "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT ?",
			inArgs:  []any{1, 2, 3, 4},
			expSQL:  "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT CAST(? AS Uint64)",
			expArgs: []any{1, 2, 3, 4},
		},
		{
			name:    "parametrized limit+offset with other args",
			inSQL:   "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT ? OFFSET ?",
			inArgs:  []any{1, 2, 3, 4},
			expSQL:  "SELECT ?, ? FROM tbl WHERE id = ? ORDER BY score LIMIT CAST(? AS Uint64) OFFSET CAST(? AS Uint64)",
			expArgs: []any{1, 2, 3, 4},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &CastLimitOffsetToUint64{}
			outSQL, outArgs, err := f.Rebind(tt.inSQL, tt.inArgs...)
			require.NoError(t, err)
			require.Equal(t, tt.expSQL, outSQL)
			require.Equal(t, tt.expArgs, outArgs)
		})
	}
}
