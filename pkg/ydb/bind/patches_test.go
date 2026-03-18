package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchesUnused(t *testing.T) {
	p := &Patches{
		patches: map[string]*patchEntry{
			"SELECT foo": {name: "001", sql: "SELECT ydb_foo"},
			"UPDATE t":   {name: "002", sql: "UPSERT t"},
		},
	}
	require.ElementsMatch(t, []string{"001", "002"}, p.UnusedPatches())

	_, _, _ = p.Rebind("SELECT foo")
	require.ElementsMatch(t, []string{"002"}, p.UnusedPatches())

	_, _, _ = p.Rebind("UPDATE t")
	require.Empty(t, p.UnusedPatches())
}
