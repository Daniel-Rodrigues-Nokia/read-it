// Package aisum is responsible to host multiplie help functions to deal with ai summary of FE tests
package aisum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func buildInstructions(tests []string) string {
	s := strings.Builder{}
	s.WriteString(Instructions)

	for _, test := range tests {
		s.WriteString("\n")
		s.WriteString(test)
		s.WriteString("\n")
	}

	return s.String()
}

func transformPayload(tests []string) (*bytes.Buffer, error) {
	message := map[string]string{"role": "user", "content": buildInstructions(tests)}

	payload := map[string]any{"model": "gemma3:12b", "messages": []map[string]string{message}}

	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(payloadInBytes), nil
}

func SummarizeTests(tests []string) (*http.Response, context.CancelFunc, error) {
	// load .env vars
	loadedVars, err := internal.LoadEnv("API_KEY", "HOST", "PORT", "ENDPOINT")
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	payloadStringify, err := transformPayload(tests)
	if err != nil {
		cancel()
		return nil, nil, err
	}

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
