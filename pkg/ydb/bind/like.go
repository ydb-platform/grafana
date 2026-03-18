package bind

import (
	"database/sql/driver"
	"strings"
)

var _ Binder = (*ConvertLikeToStartsWithEndsWith)(nil)

type ConvertLikeToStartsWithEndsWith struct{}

func findLike(sql string, start int) int {
	inString := false
	var quote byte
	i := start
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
		if (c == 'L' || c == 'l') && i+4 <= len(sql) {
			before := i == 0 || sql[i-1] == ' ' || sql[i-1] == '\t' || sql[i-1] == '\n' || sql[i-1] == '\r'
			after := i+4 == len(sql) || !isWordByte(sql[i+4])
			if before && after && strings.EqualFold(sql[i:i+4], "LIKE") {
				return i
			}
		}
		i++
	}
	return -1
}

func likeRightOperand(sql string, likeEnd int) (start, end int, isPlaceholder bool, quote byte, isConcat bool) {
	start = likeEnd
	for start < len(sql) && (sql[start] == ' ' || sql[start] == '\t' || sql[start] == '\n' || sql[start] == '\r') {
		start++
	}
	if start >= len(sql) {
		return start, start, false, 0, false
	}
	if sql[start] == '?' {
		return start, start + 1, true, 0, false
	}
	if sql[start] == '\'' || sql[start] == '"' {
		quote = sql[start]
		end = start + 1
		for end < len(sql) {
			if sql[end] == quote {
				if end+1 < len(sql) && sql[end+1] == quote {
					end += 2
					continue
				}
				end++
				literalEnd := end
				for {
					p := end
					for p < len(sql) && (sql[p] == ' ' || sql[p] == '\t' || sql[p] == '\n' || sql[p] == '\r') {
						p++
					}
					if p+2 > len(sql) || sql[p] != '|' || sql[p+1] != '|' {
						break
					}
					p += 2
					for p < len(sql) && (sql[p] == ' ' || sql[p] == '\t' || sql[p] == '\n' || sql[p] == '\r') {
						p++
					}
					if p >= len(sql) {
						break
					}
					if sql[p] == '?' {
						end = p + 1
						continue
					}
					if sql[p] == '\'' || sql[p] == '"' {
						q := sql[p]
						end = p + 1
						for end < len(sql) {
							if sql[end] == q {
								if end+1 < len(sql) && sql[end+1] == q {
									end += 2
									continue
								}
								end++
								break
							}
							if sql[end] == '\\' && end+1 < len(sql) {
								end += 2
								continue
							}
							end++
						}
						continue
					}
					break
				}
				if end > literalEnd {
					return start, end, false, quote, true
				}
				return start, literalEnd, false, quote, false
			}
			if sql[end] == '\\' && end+1 < len(sql) {
				end += 2
				continue
			}
			end++
		}
		return start, end, false, quote, false
	}
	return start, start, false, 0, false
}

func escapeLiteral(s string, quote byte) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == quote {
			b.WriteByte(quote)
			b.WriteByte(quote)
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func likePatternKind(s string) (kind string, trimmed string) {
	if len(s) == 0 {
		return "", s
	}
	last := len(s) - 1
	if s[last] == '%' {
		for i := 0; i < last; i++ {
			if s[i] == '%' {
				return "", s
			}
		}
		return "start", s[:last]
	}
	if s[0] == '%' {
		for i := 1; i < len(s); i++ {
			if s[i] == '%' {
				return "", s
			}
		}
		return "end", s[1:]
	}
	return "", s
}

func (f *ConvertLikeToStartsWithEndsWith) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	if !strings.Contains(sql, " LIKE ") && !strings.Contains(sql, " like ") {
		return sql, args, nil
	}
	var qPositions []int
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
			i++
			continue
		}
		if c == '\'' || c == '"' {
			inString = true
			quote = c
			i++
			continue
		}
		if c == '?' {
			qPositions = append(qPositions, i)
		}
	}
	argsCopy := make([]driver.NamedValue, len(args))
	for i := range args {
		argsCopy[i] = args[i]
	}
	argIdx := 0
	var out strings.Builder
	pos := 0
	for {
		likePos := findLike(sql, pos)
		if likePos < 0 {
			out.WriteString(sql[pos:])
			break
		}
		leftStart, leftEnd := leftOperand(sql, likePos)
		rightStart, rightEnd, isPlaceholder, quoteChar, isConcat := likeRightOperand(sql, likePos+4)
		out.WriteString(sql[pos:leftStart])
		var pattern string
		var argIndex int
		if isConcat {
			firstLiteralEnd := rightStart + 1
			for firstLiteralEnd < rightEnd {
				if sql[firstLiteralEnd] == quoteChar {
					if firstLiteralEnd+1 < rightEnd && sql[firstLiteralEnd+1] == quoteChar {
						firstLiteralEnd += 2
						continue
					}
					break
				}
				if sql[firstLiteralEnd] == '\\' && firstLiteralEnd+1 < rightEnd {
					firstLiteralEnd += 2
					continue
				}
				firstLiteralEnd++
			}
			if firstLiteralEnd < rightEnd && sql[firstLiteralEnd] == quoteChar {
				firstLiteralEnd++
			}
			expr := sql[rightStart:rightEnd]
			litEnd := firstLiteralEnd - rightStart
			if firstLiteralEnd > rightStart+1 && sql[rightStart+1] == '%' && litEnd >= 2 {
				out.WriteString("EndsWith(")
				out.WriteString(sql[leftStart:leftEnd])
				out.WriteString(", ")
				out.WriteString(expr[0:1])
				out.WriteString(expr[2 : litEnd-1])
				out.WriteString(expr[litEnd-1 : litEnd])
				out.WriteString(expr[litEnd:])
				out.WriteString(")")
			} else {
				out.WriteString("StartsWith(")
				out.WriteString(sql[leftStart:leftEnd])
				out.WriteString(", ")
				out.WriteString(expr)
				out.WriteString(")")
			}
			pos = rightEnd
			argIdx++
			continue
		}
		if isPlaceholder {
			for ai, qp := range qPositions {
				if qp == rightStart {
					argIndex = ai
					break
				}
			}
			if argIndex >= len(argsCopy) {
				out.WriteString(sql[leftStart:rightEnd])
				pos = rightEnd
				continue
			}
			if s, ok := argsCopy[argIndex].Value.(string); ok {
				pattern = s
			}
		} else {
			pattern = sql[rightStart+1 : rightEnd-1]
			var b strings.Builder
			for i := 0; i < len(pattern); i++ {
				if i+1 < len(pattern) && pattern[i] == '\'' && pattern[i+1] == '\'' {
					b.WriteByte('\'')
					i++
					continue
				}
				b.WriteByte(pattern[i])
			}
			pattern = b.String()
		}
		kind, trimmed := likePatternKind(pattern)
		if kind == "start" {
			out.WriteString("StartsWith(")
			out.WriteString(sql[leftStart:leftEnd])
			if isPlaceholder {
				out.WriteString(", ?)")
				argsCopy[argIndex].Value = trimmed
			} else {
				if quoteChar == 0 {
					quoteChar = '\''
				}
				out.WriteString(", ")
				out.WriteByte(quoteChar)
				out.WriteString(escapeLiteral(trimmed, quoteChar))
				out.WriteByte(quoteChar)
				out.WriteString(")")
			}
			pos = rightEnd
		} else if kind == "end" {
			out.WriteString("EndsWith(")
			out.WriteString(sql[leftStart:leftEnd])
			if isPlaceholder {
				out.WriteString(", ?)")
				argsCopy[argIndex].Value = trimmed
			} else {
				if quoteChar == 0 {
					quoteChar = '\''
				}
				out.WriteString(", ")
				out.WriteByte(quoteChar)
				out.WriteString(escapeLiteral(trimmed, quoteChar))
				out.WriteByte(quoteChar)
				out.WriteString(")")
			}
			pos = rightEnd
		} else {
			out.WriteString(sql[leftStart:rightEnd])
			pos = rightEnd
		}
		argIdx++
	}
	return out.String(), argsCopy, nil
}
