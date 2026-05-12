package reader

import "strings"

const (
	startOfTest          string = "it("
	startOfMultipleTests string = "["
)

type Cypress struct{}

func (c Cypress) IsValidLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)

	switch {
	case strings.HasPrefix(trimmedLine, startOfTest):
		return true
	case strings.HasPrefix(trimmedLine, startOfMultipleTests):
		return true
	default:
		return false
	}
}

func (c Cypress) GetTestTitle(test string) string {
	for line := range strings.SplitSeq(test, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, startOfTest) {
			// grab what's inside it("...") before the comma or closing parenthesis
			start := strings.Index(line, "(")
			end := strings.LastIndex(line, ",")
			if start != -1 && end > start {
				return line[start+1 : end]
			}
			return line
		}
	}
	return "(unknown test)"
}

func (c Cypress) GetInstructions() string {
	return ""
}
