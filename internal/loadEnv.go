package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	jiraAPIKey  string
	jiraURL     string
	jiraProject string
	jiraEmail   string
}

const env = "read-it.env"

func loadEnv() (*config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	envPath := fmt.Sprintf("%s%s%s", configDir, string(os.PathSeparator), env)

	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("failed to load %s file. Creating base template.", envPath)
		// TODO: generate template file
		return nil, nil
	}

	return &config{
		jiraAPIKey:  os.Getenv("JIRA_API_KEY"),
		jiraURL:     os.Getenv("JIRA_URL"),
		jiraProject: os.Getenv("JIRA_PROJECT"),
		jiraEmail:   os.Getenv("JIRA_EMAIL"),
	}, nil
}
