package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*DropIndex)(nil)
)

const dropIndexPrefix = "DROP INDEX "

type DropIndex struct{}

func (DropIndex) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if len(sql) < len(dropIndexPrefix) || !strings.EqualFold(sql[:len(dropIndexPrefix)], dropIndexPrefix) {
		return sql
	}
	// Parse index name and where it ends (rest of sql)
	i := len(dropIndexPrefix)
	for i < len(sql) && (sql[i] == ' ' || sql[i] == '\t') {
		i++
	}
	start := i
	if start < len(sql) && sql[start] == '`' {
		start++
		for i < len(sql) && sql[i] != '`' {
			i++
		}
		if i < len(sql) {
			i++
		}
	} else {
		for i < len(sql) && sql[i] != ' ' && sql[i] != ';' && sql[i] != '\t' {
			i++
		}
	}
	indexName := strings.Trim(sql[start:i], " `")
	rest := sql[i:]

	tableName := ""
	if table != nil {
		tableName = table.Name
	} else if dialect != nil {
		tables, err := dialect.GetTables()
		if err != nil || len(tables) == 0 {
			return sql
		}
		tableSet := make(map[string]struct{}, len(tables))
		for _, t := range tables {
			tableSet[t.Name] = struct{}{}
		}
		if len(indexName) > 4 && (indexName[:4] == "UQE_" || indexName[:4] == "uqe_") {
			suffix := indexName[4:]
			for j := 0; j < len(suffix); j++ {
				if suffix[j] == '_' {
					candidate := suffix[:j]
					if _, ok := tableSet[candidate]; ok {
						tableName = candidate
						break
					}
				}
			}
			if tableName == "" {
				if _, ok := tableSet[suffix]; ok {
					tableName = suffix
				}
			}
		}
	}
	if tableName == "" {
		return sql
	}
	// Output canonical "DROP INDEX" so case is preserved regardless of input.
	return "ALTER TABLE `" + tableName + "` DROP INDEX " + indexName + rest
}
