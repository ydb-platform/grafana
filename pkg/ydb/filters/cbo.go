package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*EnableCostBasedOptimizer)(nil)
)

type EnableCostBasedOptimizer struct{}

func (f *EnableCostBasedOptimizer) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if f.shouldUseCostBasedOptimization(sql) {
		return f.prependPragma(sql)
	}
	return sql
}

// shouldUseCostBasedOptimization returns true for "heavy" queries that typically benefit from
// ydb.CostBasedOptimization (dashboard/folder list with permission subqueries). Returns false for
// simple lookups where the pragma often produces worse plans.
func (f *EnableCostBasedOptimizer) shouldUseCostBasedOptimization(query string) bool {
	q := strings.ToUpper(query)
	hasPermission := strings.Contains(q, "PERMISSION")
	hasDashboardOrFolder := strings.Contains(q, "DASHBOARD") || strings.Contains(q, "FOLDER")
	hasOrderOrLimit := strings.Contains(q, "ORDER BY") || strings.Contains(q, "LIMIT")
	// Heavy pattern: permission + dashboard/folder + ordering/limit (dashboard search, folder list)
	if hasPermission && hasDashboardOrFolder && hasOrderOrLimit {
		return true
	}
	// Multiple nested IN (SELECT ...) — complex permission/folder chain
	inSelectCount := 0
	for i := 0; i < len(q)-10; i++ {
		if strings.HasPrefix(q[i:], "IN (SELECT") {
			inSelectCount++
			if inSelectCount >= 2 {
				return true
			}
		}
	}
	return false
}

const ydbCostBasedOptimizationPragma = `PRAGMA ydb.CostBasedOptimization = "on";` + "\n"

func (f *EnableCostBasedOptimizer) prependPragma(query string) string {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return query
	}
	if strings.HasPrefix(trimmed, strings.TrimSpace(ydbCostBasedOptimizationPragma)) {
		return query
	}
	return ydbCostBasedOptimizationPragma + query
}
