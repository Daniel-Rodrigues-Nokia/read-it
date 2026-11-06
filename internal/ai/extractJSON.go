package ai

import (
	"regexp"
	"strings"
)

var fencedJSONRE = regexp.MustCompile(`(?s)^` + "```" + `(?:json)?\s*(.*?)\s*` + "```" + `$`)

func extractJSON(content string) string {
	trimmed := strings.TrimSpace(content)
	if m := fencedJSONRE.FindStringSubmatch(trimmed); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")

	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(trimmed[start : end+1])
	}

	return strings.TrimSpace(trimmed)
}
