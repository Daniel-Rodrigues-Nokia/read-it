// Package reader
package reader

import (
	"bufio"
	"os"
	"strings"

	i "read-it/internal/item"
)

//////////////////////
//
//  Interface
//
//////////////////////

type Reader interface {
	IsValidLine(line string) bool
	GetTestTitle(test string) string
}

//////////////////////
//
// Private funcs
//
//////////////////////

func updateDepthCounter(rawLine *string, counter *int) bool {
	hasCounterBeenModified := false

	for _, char := range *rawLine {
		if char == '(' || char == '[' || char == '{' {
			*counter++
			hasCounterBeenModified = true
		}

		if char == ')' || char == ']' || char == '}' {
			*counter--
			hasCounterBeenModified = true
		}
	}

	return hasCounterBeenModified
}

//////////////////////
//
// Public API
//
//////////////////////

func ScanTests(filePath string, reader Reader) ([]i.Item, error) {
	// get file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read file
	scanner := bufio.NewScanner(file)

	counter := 0
	canScan := false
	outPut := strings.Builder{}
	tests := make([]i.Item, 0)

	for scanner.Scan() {
		rawLine := scanner.Text()

		if reader.IsValidLine(rawLine) {
			canScan = true
		}

		// skip scan otherwise
		if !canScan {
			continue
		}

		hasBeenMod := updateDepthCounter(&rawLine, &counter)

		outPut.WriteString(rawLine + "\n")

		if !hasBeenMod {
			continue
		}

		// reset state
		if counter <= 0 {
			counter = 0
			canScan = false

			// get test and push to slice
			// TODO: maybe i.NewItem can call GetTestTitle internally ??
			tests = append(tests, i.NewItem(reader.GetTestTitle(outPut.String()), outPut.String(), false))
			outPut.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}
