package ai

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

type Cursor struct{}

func (c Cursor) GetName() string {
	return "Cursor"
}

func (c Cursor) Summarize(instructions string, tests []string) (*[]Summary, error) {
	if !isCursorInstalled() {
		return nil, errors.New("cursor cli was not detected in your system. Please install it first")
	}

	prompt, err := BuildInstructions(instructions, tests)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("agent", "--output-format", "text")

	// since instructions + test cases can be reaaaaaaally big, feed .Stdin with that (instead of passing in .Command)
	cmd.Stdin = strings.NewReader(prompt.String())

	summary, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var testSum Summary
	if err := json.Unmarshal([]byte(summary), &testSum); err != nil {
		return nil, err
	}

	return &[]Summary{testSum}, nil
}

func isCursorInstalled() bool {
	cmd := exec.Command("agent", "--version")
	err := cmd.Run()

	return err == nil
}
