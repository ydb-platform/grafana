package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLower(t *testing.T) {
	for _, tt := range []struct {
		sql string
		exp string
	}{
		{
			sql: `SELECT LOWER("absence");`,
			exp: `SELECT Unicode::ToLower(CAST("absence" AS Text));`,
		},
		{
			sql: `SELECT LOWER(name) FROM tbl;`,
			exp: `SELECT Unicode::ToLower(CAST(name AS Text)) FROM tbl;`,
		},
		{
			sql: `SELECT LOWER(t.name) FROM tbl AS t;`,
			exp: `SELECT Unicode::ToLower(CAST(t.name AS Text)) FROM tbl AS t;`,
		},
		{
			sql: "SELECT LOWER(`t`.`name`) FROM `tbl` AS `t`;",
			exp: "SELECT Unicode::ToLower(CAST(`t`.`name` AS Text)) FROM `tbl` AS `t`;",
		},
		{
			sql: `SELECT LOWERED_COLUMN FROM tbl;`,
			exp: `SELECT LOWERED_COLUMN FROM tbl;`,
		},
		{
			sql: `SELECT COLUMN_LOWER FROM tbl;`,
			exp: `SELECT COLUMN_LOWER FROM tbl;`,
		},
		{
			sql: "SELECT `id`, `uid`, `version`, `email`, `name`, `login`, `password`, `salt`, `rands`, `company`, `email_verified`, `theme`, `help_flags1`, `is_disabled`, `is_admin`, `is_service_account`, `org_id`, `created`, `updated`, `last_seen_at`, `is_provisioned` FROM `user` WHERE (LOWER(email)=LOWER(?) OR LOWER(login)=LOWER(?)) LIMIT 1",
			exp: "SELECT `id`, `uid`, `version`, `email`, `name`, `login`, `password`, `salt`, `rands`, `company`, `email_verified`, `theme`, `help_flags1`, `is_disabled`, `is_admin`, `is_service_account`, `org_id`, `created`, `updated`, `last_seen_at`, `is_provisioned` FROM `user` WHERE (Unicode::ToLower(CAST(email AS Text))=Unicode::ToLower(CAST(? AS Text)) OR Unicode::ToLower(CAST(login AS Text))=Unicode::ToLower(CAST(? AS Text))) LIMIT 1",
		},
	} {
		t.Run(tt.sql, func(t *testing.T) {
			f := Lower{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
