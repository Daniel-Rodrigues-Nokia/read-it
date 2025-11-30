// Package junit
package junit

import "strings"

type JUnit struct{}

func (j JUnit) IsValidLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)

	switch {
	case strings.HasPrefix(trimmedLine, "@Test"):
		return true
	default:
		return false
	}
}

func (j JUnit) GetTestTitle(test string) string {
	for line := range strings.SplitSeq(test, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "void") {
			// grab what's after 'void' and before "()"
			start := strings.Index(line, " ")
			end := strings.LastIndex(line, "(")
			if start != -1 && end > start {
				return line[start+1 : end]
			}
			return line
		}
	}
	return "(unknown test)"
}
