// Package ai defines AI interface
package ai

import (
	"strings"
)

type AiClient interface {
	Summarize(instructions string, tests []string) (*[]string, error)
}

////////////
//
// Public API
//
////////////

func AISummary(instructions string, tests []string, ai AiClient) (*[]string, error) {
	resp, err := ai.Summarize(instructions, tests)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

////////////
//
// Helpers
//
////////////

// NOTE: DEPRECATED: since instructions.txt will be embed into the bin file
// we can just use it straight away (no need to read files)
// func ReadInstructions(filePath string) (*strings.Builder, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
//
// 	instructions := strings.Builder{}
// 	scanner := bufio.NewScanner(file)
//
// 	for scanner.Scan() {
// 		rawLine := scanner.Text()
// 		instructions.WriteString(rawLine)
// 	}
//
// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return &instructions, nil
// }

func BuildInstructions(instructions string, tests []string) (*strings.Builder, error) {
	s := strings.Builder{}

	s.WriteString(instructions)

	for _, test := range tests {
		s.WriteString("\n")
		s.WriteString(test)
		s.WriteString("\n")
	}

	return &s, nil
}
