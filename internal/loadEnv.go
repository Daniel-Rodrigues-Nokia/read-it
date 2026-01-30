package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JiraAPIKey  string
	JiraURL     string
	JiraProject string
	JiraEmail   string
}

const ENV = "read-it.env"

func generateBaseTemplate(path string) error {
	content := []byte("JIRA_API_KEY=\nJIRA_URL=\nJIRA_PROJECT=\nJIRA_EMAIL=\n")
	err := os.WriteFile(path, content, 0644) // 0644 -> chmod +x 0644, aka, owner has read+write permissions, while group & others are read only
	if err != nil {
		return err
	}

	return nil
}

func LoadEnv() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	envPath := fmt.Sprintf("%s%s%s", configDir, string(os.PathSeparator), ENV)

	err = godotenv.Load(envPath)
	if err != nil {
		err := generateBaseTemplate(envPath)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to load %s file. Creating base template in %s", ENV, envPath)
	}

	return &Config{
		JiraAPIKey:  os.Getenv("JIRA_API_KEY"),
		JiraURL:     os.Getenv("JIRA_URL"),
		JiraProject: os.Getenv("JIRA_PROJECT"),
		JiraEmail:   os.Getenv("JIRA_EMAIL"),
	}, nil
}
