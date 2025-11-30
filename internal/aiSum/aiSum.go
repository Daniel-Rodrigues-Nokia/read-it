// Package aisum is responsible to host multiplie help functions to deal with ai summary of FE tests
package aisum

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"read-it/internal"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	LogProbs     string  `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type AIResponse struct {
	ID      string   `json:"id"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func readInstructions(filePath string) (*strings.Builder, error) {
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

func buildInstructions(filePath string, tests []string) (*strings.Builder, error) {
	s := strings.Builder{}

	inst, err := readInstructions(filePath)
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

func transformPayload(filePath string, tests []string) (*bytes.Buffer, error) {
	content, err := buildInstructions(filePath, tests)
	if err != nil {
		return nil, err
	}

	message := map[string]string{"role": "user", "content": content.String()}

	payload := map[string]any{"model": "gemma3:12b", "messages": []map[string]string{message}}

	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(payloadInBytes), nil
}

func SummarizeTests(instructionsFilePath string, tests []string) (*http.Response, context.CancelFunc, error) {
	// load .env vars
	loadedVars, err := internal.LoadEnv("API_KEY", "HOST", "PORT", "ENDPOINT")
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	payloadStringify, err := transformPayload(instructionsFilePath, tests)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	// TODO: PORT might not be needed
	URL := fmt.Sprintf("http://%s:%s/%s", loadedVars[1], loadedVars[2], loadedVars[3])

	req, err := http.NewRequestWithContext(ctx, "POST", URL, payloadStringify)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", loadedVars[0]))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	return resp, cancel, nil
}

func ReadSummary(resp *http.Response) (*AIResponse, error) {
	aiResp := new(AIResponse)

	err := json.NewDecoder(resp.Body).Decode(aiResp)
	if err != nil {
		return nil, err
	}

	return aiResp, nil
}
