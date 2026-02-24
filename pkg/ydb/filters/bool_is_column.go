package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*BoolIsColumn)(nil)
)

type BoolIsColumn struct{}

func (BoolIsColumn) Do(sql string, dialect core.Dialect, table *core.Table) string {
	var out strings.Builder
	pos := 0
	i := 0
	for i < len(sql) {
		// Skip string literals
		if sql[i] == '"' || sql[i] == '\'' {
			quote := sql[i]
			out.WriteString(sql[pos:i])
			out.WriteByte(sql[i])
			i++
			pos = i
			for i < len(sql) {
				if sql[i] == '\\' && i+1 < len(sql) {
					out.WriteByte(sql[i])
					out.WriteByte(sql[i+1])
					i += 2
					continue
				}
				if sql[i] == quote {
					out.WriteByte(sql[i])
					i++
					break
				}
				out.WriteByte(sql[i])
				i++
			}
			pos = i
			continue
		}

		// Look for "is_" then identifier (column name starting with is_)
		if i+3 <= len(sql) && sql[i:i+3] == "is_" {
			// Token may be qualified: t.is_user; find start of token
			start := i
			for start > 0 && (sql[start-1] == '.' || isWordChar(sql[start-1])) {
				start--
			}
			// Extend to full identifier: is_ + word chars
			end := i + 3
			for end < len(sql) && isWordChar(sql[end]) {
				end++
			}
			if end == i+3 {
				i++
				continue
			}
			j := end
			for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
				j++
			}
			if j >= len(sql) || sql[j] != '=' {
				i++
				continue
			}
			j++
			for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
				j++
			}
			if j >= len(sql) || sql[j] != '1' {
				i++
				continue
			}
			if j+1 < len(sql) && isWordChar(sql[j+1]) {
				i++
				continue
			}
			out.WriteString(sql[pos:start])
			out.WriteString(sql[start:end])
			pos = j + 1
			i = j + 1
			continue
		}
		i++
	}
	out.WriteString(sql[pos:])
	return out.String()
}
