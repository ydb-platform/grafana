package libraryelements

import (
	"testing"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/stretchr/testify/require"
)

func TestWriteLibraryElementFolderDashboardJoin_YDBUsesEqualityOnlyInON(t *testing.T) {
	var b db.SQLBuilder
	writeLibraryElementFolderDashboardJoin(&b, migrator.NewYDBDialect())
	s := b.GetSQLString()
	require.Contains(t, s, "ON le.folder_id = dashboard.id")
	require.NotContains(t, s, "folder_id <> 0")
}

func TestWriteParamSelectorSQLAfterFolderDashboardJoin_YDB(t *testing.T) {
	var b db.SQLBuilder
	d := migrator.NewYDBDialect()
	writeParamSelectorSQLAfterFolderDashboardJoin(&b, d, Pair{"org_id", int64(1)})
	s := b.GetSQLString()
	require.Contains(t, s, "le.folder_id <> 0")
	require.Contains(t, s, "le.org_id=?")
}
