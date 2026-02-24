package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*Replace)(nil)
)

type Replace struct{}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

const replaceFrom = "REPLACE"
const replaceTo = "Unicode::ReplaceAll"

func (Replace) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if !strings.Contains(sql, "REPLACE") {
		return sql
	}
	var out strings.Builder
	i := 0
	for i < len(sql) {
		replaced := false
		if i+len(replaceFrom) <= len(sql) && sql[i:i+len(replaceFrom)] == replaceFrom {
			prevOK := i == 0 || !isWordChar(sql[i-1])
			nextOK := i+len(replaceFrom) >= len(sql) || !isWordChar(sql[i+len(replaceFrom)])
			openParen := i + len(replaceFrom)
			for openParen < len(sql) && (sql[openParen] == ' ' || sql[openParen] == '\t') {
				openParen++
			}
			if prevOK && nextOK && openParen < len(sql) && sql[openParen] == '(' {
				// REPLACE(arg1, arg2, arg3) → Unicode::ReplaceAll(CAST(arg1 AS Text), arg2, arg3)
				// Запятые/скобки внутри строковых литералов ('...', "...", `...`) не учитываем.
				start := openParen + 1
				depth := 1
				pos := start
				var commaPos int
				var closePos int
				foundComma := false
				inString := byte(0) // 0 = нет, '\'' '"' '`' = внутри такого литерала
			scanArgs:
				for pos < len(sql) {
					b := sql[pos]
					if inString != 0 {
						if b == inString {
							if inString == '\'' && pos+1 < len(sql) && sql[pos+1] == '\'' {
								pos++ // '' — экранированная кавычка в '-строке
							} else {
								inString = 0
							}
						}
						pos++
						continue
					}
					switch b {
					case '\'', '"', '`':
						inString = b
						pos++
						continue
					case '(':
						depth++
					case ')':
						depth--
						if depth == 0 {
							closePos = pos
							break scanArgs
						}
					case ',':
						if depth == 1 && !foundComma {
							commaPos = pos
							foundComma = true
						}
					}
					pos++
				}
				if foundComma && depth == 0 {
					arg1 := sql[start:commaPos]
					out.WriteString(replaceTo)
					out.WriteString(sql[i+len(replaceFrom):openParen]) // preserve space(s) after REPLACE
					out.WriteString("(CAST(")
					out.WriteString(arg1)
					out.WriteString(" AS Text)")
					out.WriteString(sql[commaPos:closePos+1])
					i = closePos + 1
					replaced = true
				}
			}
			if !replaced && prevOK && nextOK {
				out.WriteString(replaceTo)
				i += len(replaceFrom)
				continue
			}
		}
		if replaced {
			continue
		}
		out.WriteByte(sql[i])
		i++
	}
	return out.String()
}
