package bind

func isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

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

func findMatchingParen(sql string, open int) int {
	if open >= len(sql) || sql[open] != '(' {
		return -1
	}
	depth := 1
	inString := false
	var quote byte
	i := open + 1
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

func findMatchingParenCreateIndex(s string, open int) int {
	return findMatchingParen(s, open)
}

func findIgnoreCase(sql, s string) int {
	if len(s) == 0 || len(sql) < len(s) {
		return -1
	}
	for i := 0; i <= len(sql)-len(s); i++ {
		var j int
		for j = 0; j < len(s); j++ {
			c, d := sql[i+j], s[j]
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			if d >= 'A' && d <= 'Z' {
				d += 'a' - 'A'
			}
			if c != d {
				break
			}
		}
		if j == len(s) && (i == 0 || !isWordByte(sql[i-1])) {
			return i
		}
	}
	return -1
}

func leftOperand(sql string, likePos int) (start, end int) {
	end = likePos
	for end > 0 && (sql[end-1] == ' ' || sql[end-1] == '\t' || sql[end-1] == '\n' || sql[end-1] == '\r') {
		end--
	}
	if end == 0 {
		return 0, likePos
	}
	start = end
	if sql[end-1] == '?' {
		start = end - 1
		return start, end
	}
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
