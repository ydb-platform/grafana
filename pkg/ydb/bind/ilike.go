package bind

import (
	"database/sql/driver"
	"strings"
)

var _ Binder = (*ConvertILikeToLikeLowerCase)(nil)

type ConvertILikeToLikeLowerCase struct{}

func findILike(sql string, start int) int {
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
		if (c == 'I' || c == 'i') && i+5 <= len(sql) {
			if strings.EqualFold(sql[i:i+5], "ILIKE") && (i == 0 || !isWordByte(sql[i-1])) {
				return i
			}
		}
		i++
	}
	return -1
}

func rightOperand(sql string, ilikeEnd int) (start, end int) {
	start = ilikeEnd
	for start < len(sql) && (sql[start] == ' ' || sql[start] == '\t' || sql[start] == '\n' || sql[start] == '\r') {
		start++
	}
	if start >= len(sql) {
		return start, start
	}
	end = start
	if sql[start] == '?' {
		end = start + 1
		return start, end
	}
	if sql[start] == '$' {
		end = start + 1
		for end < len(sql) && (isWordByte(sql[end]) || (sql[end] >= '0' && sql[end] <= '9')) {
			end++
		}
		return start, end
	}
	for end < len(sql) && (isWordByte(sql[end]) || sql[end] == '.' || sql[end] == '`') {
		if sql[end] == '`' {
			end++
			for end < len(sql) && sql[end] != '`' {
				end++
			}
			if end < len(sql) {
				end++
			}
			continue
		}
		end++
	}
	return start, end
}

func (f *ConvertILikeToLikeLowerCase) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	if !strings.Contains(sql, "ILIKE") && !strings.Contains(sql, "ilike") {
		return sql, args, nil
	}
	var out strings.Builder
	pos := 0
	for {
		ilikePos := findILike(sql, pos)
		if ilikePos < 0 {
			out.WriteString(sql[pos:])
			break
		}
		leftStart, leftEnd := leftOperand(sql, ilikePos)
		rightStart, rightEnd := rightOperand(sql, ilikePos+5)
		out.WriteString(sql[pos:leftStart])
		out.WriteString("Unicode::ToLower(")
		out.WriteString(sql[leftStart:leftEnd])
		out.WriteString(") LIKE Unicode::ToLower(")
		out.WriteString(sql[rightStart:rightEnd])
		out.WriteString(")")
		pos = rightEnd
	}
	return out.String(), args, nil
}
