package bind

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

var _ Binder = (*ConvertPositionalArgsToYdbNamedParameters)(nil)

type ConvertPositionalArgsToYdbNamedParameters struct{}

func (f *ConvertPositionalArgsToYdbNamedParameters) Rebind(query string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	var buf strings.Builder
	var beginQuote rune
	argIndex := 0
	paramIndex := 0
	outArgs := make([]driver.NamedValue, 0, len(args))
	for _, c := range query {
		if beginQuote == 0 {
			switch c {
			case '\'', '"':
				beginQuote = c
				buf.WriteRune(c)
			case '?':
				if argIndex < len(args) {
					val := args[argIndex].Value
					if val == nil {
						buf.WriteString("NULL")
					} else {
						paramIndex++
						paramName := fmt.Sprintf("p%d", paramIndex)
						buf.WriteString("$" + paramName)
						outArgs = append(outArgs, driver.NamedValue{
							Name:  paramName,
							Value: val,
						})
					}
				} else {
					paramIndex++
					paramName := fmt.Sprintf("p%d", paramIndex)
					buf.WriteString("$" + paramName)
					outArgs = append(outArgs, args[argIndex])
				}
				argIndex++
			default:
				buf.WriteRune(c)
			}
		} else {
			switch c {
			case beginQuote:
				beginQuote = 0
				buf.WriteRune(c)
			default:
				buf.WriteRune(c)
			}
		}
	}
	return buf.String(), outArgs, nil
}
