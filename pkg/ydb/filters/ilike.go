package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*ILike)(nil)
)

type ILike struct{}

// findILike finds the next "ILIKE" outside string literals starting at start.
// Returns position of 'I' or -1. Checks that ILIKE is not part of a longer word.
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

// leftOperand finds the left operand of ILIKE at ilikePos: walk backward from ilikePos-1.
// Returns (start, end) of the operand in sql (end is exclusive), so sql[start:end] is the operand.
func leftOperand(sql string, ilikePos int) (start, end int) {
	end = ilikePos
	for end > 0 && (sql[end-1] == ' ' || sql[end-1] == '\t' || sql[end-1] == '\n' || sql[end-1] == '\r') {
		end--
	}
	if end == 0 {
		return 0, ilikePos
	}
	start = end
	if sql[end-1] == '?' {
		start = end - 1
		return start, end
	}
	// Identifier: backticks, letters, digits, _, .
	for start > 0 {
		c := sql[start-1]
		if c == '`' {
			start--
			for start > 0 && sql[start-1] != '`' {
				start--
			}
			if start > 0 {
				start--
			}
			continue
		}
		if isWordByte(c) || c == '.' {
			start--
			continue
		}
		break
	}
	return start, end
}

// rightOperand finds the right operand of ILIKE: ilikeEnd = ilikePos+5, skip space, then take ? or $param.
// Returns (start, end) of the operand.
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
	// identifier (e.g. column name on right side - rare but handle)
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

func (f *ILike) Do(sql string, _ core.Dialect, _ *core.Table) string {
	if !strings.Contains(sql, "ILIKE") && !strings.Contains(sql, "ilike") {
		return sql
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
	return out.String()
}
