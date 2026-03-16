package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var (
	_ core.Filter         = (*ConvertLimitArgToUint64)(nil)
	_ core.FilterWithArgs = (*ConvertLimitArgToUint64)(nil)
)

type ConvertLimitArgToUint64 struct{}

func (f *ConvertLimitArgToUint64) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

// indexOfLimitPlaceholder returns the byte position of the '?' that follows "LIMIT" in sql, or -1.
// Only considers LIMIT outside string literals. LIMIT may be followed by optional whitespace then ?.
func (f *ConvertLimitArgToUint64) indexOfLimitPlaceholder(sql string) int {
	inString := false
	var quote byte
	i := 0
	for i < len(sql) {
		c := sql[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(sql) {
				i++
			}
			i++
			continue
		}
		if c == '\'' || c == '"' {
			inString = true
			quote = c
			i++
			continue
		}
		if (c == 'L' || c == 'l') && i+5 <= len(sql) {
			word := sql[i : i+5]
			if (word == "LIMIT" || word == "limit") && (i == 0 || !f.isWordByte(sql[i-1])) {
				j := i + 5
				for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t' || sql[j] == '\n' || sql[j] == '\r') {
					j++
				}
				if j < len(sql) && sql[j] == '?' {
					return j
				}
			}
		}
		i++
	}
	return -1
}

func (f *ConvertLimitArgToUint64) isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func (f *ConvertLimitArgToUint64) toUint64(v any) (uint64, bool) {
	switch x := v.(type) {
	case uint64:
		return x, true
	case uint:
		return uint64(x), true
	case uint32:
		return uint64(x), true
	case uint16:
		return uint64(x), true
	case uint8:
		return uint64(x), true
	case int:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int64:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int32:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int16:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int8:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	default:
		return 0, false
	}
}

func (f *ConvertLimitArgToUint64) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	limitQPos := f.indexOfLimitPlaceholder(sql)
	if limitQPos < 0 || len(args) == 0 {
		return sql, args
	}
	var qIndex int
	inString := false
	var quote byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(sql) {
				i++
			}
			continue
		}
		if c == '\'' || c == '"' {
			inString = true
			quote = c
			continue
		}
		if c == '?' {
			if i == limitQPos {
				break
			}
			qIndex++
		}
	}
	if qIndex >= len(args) {
		return sql, args
	}
	if u, ok := f.toUint64(args[qIndex]); ok {
		out := make([]any, len(args))
		copy(out, args)
		out[qIndex] = u
		return sql, out
	}
	return sql, args
}
