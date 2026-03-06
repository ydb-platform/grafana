package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*AddColumn)(nil)
)

const addColumnLiteral = "ADD COLUMN"

type AddColumn struct{}

func (AddColumn) Do(sql string, _ core.Dialect, table *core.Table) string {
	if table == nil {
		return sql
	}
	pos := findIgnoreCase(sql, addColumnLiteral)
	if pos < 0 {
		return sql
	}
	// column name starts after "ADD COLUMN" and optional whitespace
	start := pos + len(addColumnLiteral)
	for start < len(sql) && (sql[start] == ' ' || sql[start] == '\t') {
		start++
	}
	if start >= len(sql) {
		return sql
	}
	var name string
	if sql[start] == '`' {
		end := start + 1
		for end < len(sql) && sql[end] != '`' {
			end++
		}
		if end < len(sql) {
			name = sql[start+1 : end]
		}
	} else {
		end := start
		for end < len(sql) && isWordByte(sql[end]) {
			end++
		}
		name = sql[start:end]
	}
	if name == "" {
		return sql
	}
	for _, col := range table.Columns() {
		if col.Name == name {
			return "SELECT * FROM (SELECT 1) WHERE false; -- nop"
		}
	}
	return sql
}

// findIgnoreCase returns index of first occurrence of s in sql (case-insensitive), or -1.
func findIgnoreCase(sql, s string) int {
	if len(s) == 0 || len(sql) < len(s) {
		return -1
	}
	for i := 0; i <= len(sql)-len(s); i++ {
		var j int
		for j = 0; j < len(s); j++ {
			c, d := sql[i+j], s[j]
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			if d >= 'A' && d <= 'Z' {
				d += 'a' - 'A'
			}
			if c != d {
				break
			}
		}
		if j == len(s) && (i == 0 || !isWordByte(sql[i-1])) {
			return i
		}
	}
	return -1
}
