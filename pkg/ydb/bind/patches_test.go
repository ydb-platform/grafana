package bind

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_minify(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out string
	}{
		{
			in:  "SELECT `u`.`id`, `u`.`uid`, `u`.`email`, `u`.`name`, `u`.`login`, `u`.`is_admin`, `u`.`is_disabled`, `u`.`last_seen_at`, `user_auth`.`auth_module`, `u`.`is_provisioned` FROM `user` AS `u` LEFT JOIN user_auth ON user_auth.id=(\n\t\tSELECT id from user_auth\n\t\t\tWHERE user_auth.user_id = u.id\n\t\t\tORDER BY user_auth.created DESC  LIMIT 1) WHERE (u.is_service_account = ? AND  1 = 1) AND (last_seen_at > ?) ORDER BY `u`.`login` ASC, `u`.`email` ASC LIMIT 50",
			out: "SELECT `u`.`id`, `u`.`uid`, `u`.`email`, `u`.`name`, `u`.`login`, `u`.`is_admin`, `u`.`is_disabled`, `u`.`last_seen_at`, `user_auth`.`auth_module`, `u`.`is_provisioned` FROM `user` AS `u` LEFT JOIN user_auth ON user_auth.id=( SELECT id from user_auth WHERE user_auth.user_id = u.id ORDER BY user_auth.created DESC  LIMIT 1) WHERE (u.is_service_account = ? AND  1 = 1) ORDER BY `u`.`login` ASC, `u`.`email` ASC LIMIT 50",
		},
	} {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(t, tt.out, minify(tt.in))
		})
	}
}
