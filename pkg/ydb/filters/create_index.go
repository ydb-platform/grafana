package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*CreateIndex)(nil)
)

type CreateIndex struct{}

// Do rewrites CREATE [UNIQUE] INDEX name ON table (cols) to
// ALTER TABLE table CREATE [UNIQUE] INDEX name GLOBAL ON (cols).
func (CreateIndex) Do(sql string, _ core.Dialect, _ *core.Table) string {
	s := skipSpace(sql, 0)
	if !hasPrefixIgnoreCase(s, "CREATE ") {
		return sql
	}
	s = skipSpace(s, 7)
	unique := hasPrefixIgnoreCase(s, "UNIQUE ")
	if unique {
		s = skipSpace(s, 7)
	}
	if !hasPrefixIgnoreCase(s, "INDEX ") {
		return sql
	}
	s = skipSpace(s, 6)
	indexName, indexQuoted, s := parseIdentifier(s)
	if indexName == "" {
		return sql
	}
	s = skipSpace(s, 0)
	if !hasPrefixIgnoreCase(s, "ON ") {
		return sql
	}
	s = skipSpace(s, 3)
	tableName, tableQuoted, s := parseIdentifier(s)
	if tableName == "" {
		return sql
	}
	s = skipSpace(s, 0)
	if len(s) == 0 || s[0] != '(' {
		return sql
	}
	// column list is from '(' to matching ')'
	colEnd := findMatchingParenCreateIndex(s, 0)
	if colEnd < 0 {
		return sql
	}
	colList := s[:colEnd+1]
	rest := s[colEnd+1:]

	quote := func(name string, quoted bool) string {
		if quoted {
			return "`" + name + "`"
		}
		return name
	}
	out := "ALTER TABLE " + quote(tableName, tableQuoted) + " CREATE "
	if unique {
		out += "UNIQUE "
	}
	out += "INDEX " + quote(indexName, indexQuoted) + " GLOBAL ON " + colList + rest
	return out
}

// skipSpace returns s[n:] with leading spaces/tabs/newlines removed.
func skipSpace(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	s = s[n:]
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	return s
}

func hasPrefixIgnoreCase(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		a, b := s[i], prefix[i]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}

// parseIdentifier returns (name, wasQuoted, rest). Handles `name` or bare word.
func parseIdentifier(s string) (name string, quoted bool, rest string) {
	if len(s) == 0 {
		return "", false, s
	}
	if s[0] == '`' {
		end := 1
		for end < len(s) && s[end] != '`' {
			end++
		}
		if end < len(s) {
			return s[1:end], true, s[end+1:]
		}
		return "", false, s
	}
	end := 0
	for end < len(s) && isWordByte(s[end]) {
		end++
	}
	if end > 0 {
		return s[:end], false, s[end:]
	}
	return "", false, s
}

// findMatchingParenCreateIndex finds matching ')' for s[0]=='(', respecting strings. Returns index of ')' or -1.
func findMatchingParenCreateIndex(s string, open int) int {
	if open >= len(s) || s[open] != '(' {
		return -1
	}
	depth := 1
	inString := false
	var quote byte
	i := open + 1
	for i < len(s) {
		c := s[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(s) {
				i++
			}
			i++
			continue
		}
		if c == '\'' || c == '"' || c == '`' {
			inString = true
			quote = c
			i++
			continue
		}
		if c == '(' {
			depth++
			i++
			continue
		}
		if c == ')' {
			depth--
			if depth == 0 {
				return i
			}
			i++
			continue
		}
		i++
	}
	return -1
}
