package bind

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"unicode"
)

var _ Binder = (*ConvertPositionalArgsToYdbNamedParameters)(nil)

type ConvertPositionalArgsToYdbNamedParameters struct{}

// placeholder описывает позицию параметра в запросе: либо ?, либо $name.
type placeholder struct {
	start int
	end   int
	name  string // для $name — имя; для ? — пустое
}

// unwrapArg извлекает значение из аргумента, убирая обёртки sql.NamedArg и driver.NamedValue.
func unwrapArg(arg any) any {
	switch a := arg.(type) {
	case driver.NamedValue:
		return a.Value
	case sql.NamedArg:
		return a.Value
	default:
		return arg
	}
}

func (f *ConvertPositionalArgsToYdbNamedParameters) Rebind(query string, args ...any) (string, []any, error) {
	placeholders := parsePlaceholders(query)
	usedNames := make(map[string]struct{})
	for _, p := range placeholders {
		if p.name != "" {
			usedNames[p.name] = struct{}{}
		}
	}
	paramCounter := 0
	nextName := func() string {
		for {
			paramCounter++
			name := fmt.Sprintf("p%d", paramCounter)
			if _, ok := usedNames[name]; !ok {
				return name
			}
		}
	}

	outArgs := make([]any, 0, len(placeholders))
	var buf strings.Builder
	argIndex := 0
	pos := 0
	for _, p := range placeholders {
		if argIndex >= len(args) {
			return "", nil, fmt.Errorf("not enough args for placeholders")
		}
		arg := args[argIndex]
		argIndex++
		value := unwrapArg(arg)
		var name string
		if p.name != "" {
			name = p.name
		} else if value != nil {
			// для nil не вызываем nextName(), чтобы не занимать номер (нумерация p1,p2,p3,... без пропусков для выводимых параметров)
			name = nextName()
		}
		buf.WriteString(query[pos:p.start])
		if value == nil {
			buf.WriteString("NULL")
		} else {
			outArgs = append(outArgs, driver.NamedValue{Name: name, Value: value})
			buf.WriteString("$" + name)
		}
		pos = p.end
	}
	buf.WriteString(query[pos:])
	return buf.String(), outArgs, nil
}

// parsePlaceholders возвращает список плейссхолдеров ? и $id в порядке появления (вне кавычек).
func parsePlaceholders(query string) []placeholder {
	var result []placeholder
	i := 0
	for i < len(query) {
		c := query[i]
		if c == '\'' || c == '"' {
			quote := c
			i++
			for i < len(query) {
				if query[i] == quote {
					if quote == '\'' && i+1 < len(query) && query[i+1] == '\'' {
						i += 2
						continue
					}
					i++
					break
				}
				if query[i] == '\\' && i+1 < len(query) {
					i += 2
					continue
				}
				i++
			}
			continue
		}
		if c == '?' {
			result = append(result, placeholder{start: i, end: i + 1, name: ""})
			i++
			continue
		}
		if c == '$' && i+1 < len(query) && (unicode.IsLetter(rune(query[i+1])) || query[i+1] == '_') {
			start := i
			i++
			for i < len(query) && (unicode.IsLetter(rune(query[i])) || unicode.IsDigit(rune(query[i])) || query[i] == '_') {
				i++
			}
			name := query[start+1 : i]
			result = append(result, placeholder{start: start, end: i, name: name})
			continue
		}
		i++
	}
	return result
}
