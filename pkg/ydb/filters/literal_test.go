package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLiteral(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	for _, tt := range []struct {
		name string
		in   args
		out  args
	}{
		{
			name: "IN with text literal values (single quotas)",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN ('SignUpStarted', 'InvitePending')",
				args: []any{10},
			},
			out: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (?, ?)",
				args: []any{10, "SignUpStarted", "InvitePending"},
			},
		},
		{
			name: "IN with text literal values (double quotas)",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (\"SignUpStarted\", \"InvitePending\")",
				args: []any{10},
			},
			out: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (?, ?)",
				args: []any{10, "SignUpStarted", "InvitePending"},
			},
		},
		{
			name: "IN with text literal values (backticks)",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (`SignUpStarted`, `InvitePending`)",
				args: []any{10},
			},
			out: args{
				// Backticks are not extracted (identifiers); IN filter does not apply to literal list.
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (`SignUpStarted`, `InvitePending`)",
				args: []any{10},
			},
		},
		{
			name: "IN with numeric literal values",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (1,2,3,4)",
				args: []any{10},
			},
			out: args{
				// Numeric literals are not extracted (LIMIT/OFFSET etc. handled by other filters).
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN (1,2,3,4)",
				args: []any{10},
			},
		},
		{
			name: "literal in WHERE condition",
			in: args{
				sql:  "UPDATE `temp_user` SET created = ?, updated = ? WHERE created = '0' AND status in (?,?)",
				args: []any{1772043198, 1772043198, "SignUpStarted", "InvitePending"},
			},
			out: args{
				sql:  "UPDATE `temp_user` SET created = ?, updated = ? WHERE created = ? AND status in (?,?)",
				args: []any{1772043198, 1772043198, "0", "SignUpStarted", "InvitePending"},
			},
		},
		{
			name: "literal in CREATE TABLE",
			in: args{
				sql:  "CREATE TABLE IF NOT EXISTS `alert_definition` (\n\t`id` Serial,\n\t`org_id` Int64,\n\t`title` String,\n\t`condition` String,\n\t`data` String,\n\t`updated` String,\n\t`interval_seconds` Int64 DEFAULT 60 ,\n\t`version` Int64 DEFAULT 0 ,\n\t`uid` String DEFAULT \"0\" ,\n\tPRIMARY KEY (`id`),\n\tINDEX `org_id_title` GLOBAL SYNC ON (`org_id`,`title`),\n\tINDEX `org_id_uid` GLOBAL SYNC ON (`org_id`,`uid`)\n);",
				args: nil,
			},
			out: args{
				sql:  "CREATE TABLE IF NOT EXISTS `alert_definition` (\n\t`id` Serial,\n\t`org_id` Int64,\n\t`title` String,\n\t`condition` String,\n\t`data` String,\n\t`updated` String,\n\t`interval_seconds` Int64 DEFAULT 60 ,\n\t`version` Int64 DEFAULT 0 ,\n\t`uid` String DEFAULT \"0\" ,\n\tPRIMARY KEY (`id`),\n\tINDEX `org_id_title` GLOBAL SYNC ON (`org_id`,`title`),\n\tINDEX `org_id_uid` GLOBAL SYNC ON (`org_id`,`uid`)\n);",
				args: nil,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &Literal{}
			sql, args := f.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
