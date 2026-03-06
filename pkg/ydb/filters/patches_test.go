package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatches_Do(t *testing.T) {
	f := &Patches{m: map[string]*patchEntry{
		"SELECT foo FROM bar": {
			name: "001",
			sql:  "SELECT ydb_foo FROM bar",
		},
		"UPDATE t SET x = 1": {
			name: "002",
			sql:  "UPSERT INTO t (x) VALUES (1)",
		},
	}}

	for _, tt := range []struct {
		name string
		sql  string
		exp  string
	}{
		{
			name: "lookup",
			sql:  "SELECT foo FROM bar",
			exp:  "SELECT ydb_foo FROM bar",
		},
		{
			name: "no matches",
			sql:  "SELECT other FROM table",
			exp:  "SELECT other FROM table",
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}

func TestPatchesUnused(t *testing.T) {
	p := &Patches{m: map[string]*patchEntry{
		"SELECT foo": {name: "001", sql: "SELECT ydb_foo"},
		"UPDATE t":   {name: "002", sql: "UPSERT t"},
	}}
	require.ElementsMatch(t, []string{"001", "002"}, p.UnusedPatches())

	_ = p.Do("SELECT foo", nil, nil)
	require.ElementsMatch(t, []string{"002"}, p.UnusedPatches())

	_ = p.Do("UPDATE t", nil, nil)
	require.Empty(t, p.UnusedPatches())
}
