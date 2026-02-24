package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var (
	_ core.Filter = (*RenameColumn)(nil)
)

type RenameColumn struct{}

// Do rewrites RENAME COLUMN into ADD COLUMN + UPDATE + DROP COLUMN for YDB compatibility.
// Input:  ALTER TABLE table RENAME COLUMN old_col TO new_col;
// Output: ALTER TABLE table ADD COLUMN new_col Int64;
//         UPDATE table SET new_col = old_col;
//         ALTER TABLE table DROP COLUMN old_col;
func (RenameColumn) Do(sql string, _ core.Dialect, _ *core.Table) string {
	s := skipSpace(sql, 0)
	if !hasPrefixIgnoreCase(s, "ALTER TABLE ") {
		return sql
	}
	s = skipSpace(s, 12)
	tableName, _, s := parseIdentifier(s)
	if tableName == "" {
		return sql
	}
	s = skipSpace(s, 0)
	if !hasPrefixIgnoreCase(s, "RENAME COLUMN ") {
		return sql
	}
	s = skipSpace(s, 13)
	oldCol, _, s := parseIdentifier(s)
	if oldCol == "" {
		return sql
	}
	s = skipSpace(s, 0)
	if !hasPrefixIgnoreCase(s, "TO ") {
		return sql
	}
	s = skipSpace(s, 3)
	newCol, _, _ := parseIdentifier(s)
	if newCol == "" {
		return sql
	}
	// Emit canonical uppercase statements.
	return "ALTER TABLE " + tableName + " ADD COLUMN " + newCol + " Int64;\n" +
		"UPDATE " + tableName + " SET " + newCol + " = " + oldCol + ";\n" +
		"ALTER TABLE " + tableName + " DROP COLUMN " + oldCol + ";\n"
}
