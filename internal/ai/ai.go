// Package ai defines AI interface
package ai

import (
	"bufio"
	"os"
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

func AISummary(instructionsPath string, tests []string, ai AiClient) (*[]string, error) {
	resp, err := ai.Summarize(instructionsPath, tests)
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

func ReadInstructions(filePath string) (*strings.Builder, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	instructions := strings.Builder{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		rawLine := scanner.Text()
		instructions.WriteString(rawLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &instructions, nil
}

func BuildInstructions(filePath string, tests []string) (*strings.Builder, error) {
	s := strings.Builder{}

	inst, err := ReadInstructions(filePath)
	if err != nil {
		return nil, err
	}

	s.WriteString(inst.String())

	for _, test := range tests {
		s.WriteString("\n")
		s.WriteString(test)
		s.WriteString("\n")
	}

	return &s, nil
}
