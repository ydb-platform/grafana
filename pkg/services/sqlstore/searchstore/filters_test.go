package searchstore

import (
	"testing"

	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/stretchr/testify/require"
)

func TestFolderWithAlertsFilter_Where(t *testing.T) {
	t.Run("default is correlated EXISTS", func(t *testing.T) {
		f := FolderWithAlertsFilter{}
		sql, args := f.Where()
		require.Contains(t, sql, "EXISTS")
		require.Contains(t, sql, "namespace_uid = dashboard.uid")
		require.Empty(t, args)
	})
	t.Run("YDB with org uses uncorrelated IN", func(t *testing.T) {
		f := FolderWithAlertsFilter{
			Dialect: migrator.NewYDBDialect(),
			OrgID:   99,
		}
		sql, args := f.Where()
		require.Contains(t, sql, "IN (SELECT namespace_uid FROM alert_rule")
		require.Contains(t, sql, "org_id = ?")
		require.Equal(t, []interface{}{int64(99)}, args)
	})
	t.Run("YDB with org 0 falls back to EXISTS", func(t *testing.T) {
		f := FolderWithAlertsFilter{Dialect: migrator.NewYDBDialect(), OrgID: 0}
		sql, _ := f.Where()
		require.Contains(t, sql, "EXISTS")
	})
}
