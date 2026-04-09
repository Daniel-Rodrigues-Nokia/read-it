package ai

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"

	copilot "github.com/github/copilot-sdk/go"
)

type Copilot struct {
	Model   string
	Timeout int
}

const defaultModel = "gpt-4.1"

func (cp Copilot) GetName() string {
	return "Copilot"
}

func (cp Copilot) Summarize(instructions string, tests []string) (*[]Summary, error) {
	if !isCopilotInstalled() {
		return nil, errors.New("copilot cli was not detected in your system. Please install it first")
	}

	client := copilot.NewClient(nil)
	if err := client.Start(context.Background()); err != nil {
		return nil, err
	}
	defer client.Stop()

	// set default value for Model
	if cp.Model == "" {
		cp.Model = defaultModel
	}

	session, err := client.CreateSession(context.Background(), &copilot.SessionConfig{
		Model:               cp.Model,
		OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
	})
	if err != nil {
		return nil, err
	}

	promp, err := BuildInstructions(instructions, tests)
	if err != nil {
		return nil, err
	}

	response, err := session.SendAndWait(context.Background(), copilot.MessageOptions{Prompt: promp.String()})
	if err != nil {
		return nil, err
	}

	summary := Summary{}
	err = json.Unmarshal([]byte(*response.Data.Content), &summary)
	if err != nil {
		return nil, err
	}

	return &[]Summary{summary}, nil
}

func isCopilotInstalled() bool {
	cmd := exec.Command("copilot", "--version")
	err := cmd.Run()

	return err == nil
}
