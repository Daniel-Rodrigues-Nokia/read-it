package reader

import (
	_ "embed"
	"strings"
)

//go:embed instructions-junit.txt
var instructionsForJUnit string

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
	keyword := "void " // whitespace at the end is necessary

	for line := range strings.SplitSeq(test, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, strings.TrimSpace(keyword)) || strings.HasPrefix(line, "public") {
			// grab what's after keyword and before "()"
			start := strings.Index(line, keyword)
			end := strings.LastIndex(line, "(")

			if start != -1 && end > start {
				return line[start+len(keyword) : end]
			}
			return line
		}
	}
	return "(unknown test)"
}

func (j JUnit) GetInstructions() string {
	return instructionsForJUnit
}
