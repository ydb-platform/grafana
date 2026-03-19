package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertInArgsToList(t *testing.T) {
	for _, tt := range []struct {
		name   string
		inSQL  string
		inArgs []any
		expSQL string
		expArgs []any
	}{
		{
			name:   "single IN with separate spaces",
			inSQL:  "SELECT * FROM tbl WHERE id IN (?, ?, ?)",
			inArgs: []any{1, 2, 3},
			expSQL: "SELECT * FROM tbl WHERE id IN ?",
			expArgs: []any{[]any{1, 2, 3}},
		},
		{
			name:   "IN in lower case",
			inSQL:  "SELECT * FROM tabl WHERE created > ? AND status in (?, ?)",
			inArgs: []any{1, 2, 3},
			expSQL: "SELECT * FROM tabl WHERE created > ? AND status IN ?",
			expArgs: []any{1, []any{2, 3}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertInArgsToList{}
			outSQL, outArgs, err := f.Rebind(tt.inSQL, tt.inArgs...)
			require.NoError(t, err)
			require.Equal(t, tt.expSQL, outSQL)
			require.Equal(t, tt.expArgs, outArgs)
		})
	}
}
