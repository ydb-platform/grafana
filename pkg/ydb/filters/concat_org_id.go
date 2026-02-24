package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*ConcatOrgID)(nil)
)

const castPrefix = "CAST("
const orgIdAs = "org_id AS "

type ConcatOrgID struct{}

func (ConcatOrgID) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if !strings.Contains(sql, "CAST(") || !strings.Contains(sql, "org_id AS") {
		return sql
	}
	var out strings.Builder
	pos := 0
	i := 0
	for i < len(sql) {
		if sql[i] == '"' || sql[i] == '\'' || sql[i] == '`' {
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

		if i+len(castPrefix) <= len(sql) && sql[i:i+len(castPrefix)] == castPrefix {
			castStart := i
			j := i + len(castPrefix)
			for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
				j++
			}
			colStart := j
			// Parse optional (identifier.)* then "org_id"
			for {
				if j+len(orgIdAs) <= len(sql) && sql[j:j+len(orgIdAs)] == orgIdAs {
					colEnd := j + 6 // end of "org_id"
					typeStart := j + len(orgIdAs)
					typeEnd := typeStart
					for typeEnd < len(sql) && (isWordChar(sql[typeEnd]) || sql[typeEnd] == '_') {
						typeEnd++
					}
					if typeEnd == typeStart || typeEnd >= len(sql) || sql[typeEnd] != ')' {
						i++
						break
					}
					castEnd := typeEnd
					out.WriteString(sql[pos:castStart])
					out.WriteString("Unwrap(CAST(")
					out.WriteString(sql[colStart:colEnd])
					out.WriteString(" AS STRING))")
					pos = castEnd + 1
					i = castEnd + 1
					break
				}
				// Consume identifier
				for j < len(sql) && isWordChar(sql[j]) {
					j++
				}
				if j >= len(sql) || sql[j] != '.' {
					i++
					break
				}
				j++ // consume '.'
			}
			continue
		}
		i++
	}
	out.WriteString(sql[pos:])
	return out.String()
}
