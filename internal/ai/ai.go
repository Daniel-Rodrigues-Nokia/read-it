// Package ai defines AI interface
package ai

import (
	"fmt"
	"strings"

	in "read-it/internal"
	sp "read-it/internal/components/spinner"
)

type Summary struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type AiClient interface {
	GetName() string
	Summarize(instructions string, tests []string) (*[]Summary, error)
}

////////////
//
// Public API
//
////////////

func GetAISummaryWithUI(instructions string, item []string, ag AiClient) (*[]Summary, error) {
	spinner, err := sp.NewSpinner(fmt.Sprintf("Generating summary with %s...", ag.GetName()), in.CancelCtrl).Start()
	if err != nil {
		return nil, err
	}

	defer func() {
		spinner.Stop()
		in.ClearStdOut()
	}()

	resp, err := aiSummary(instructions, item, ag)
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

func aiSummary(instructions string, tests []string, ai AiClient) (*[]Summary, error) {
	resp, err := ai.Summarize(instructions, tests)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

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
