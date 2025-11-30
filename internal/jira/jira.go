// Package jira
package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"read-it/internal"
)

//////////////////////
//
// Structs
//
//////////////////////

type Jira struct {
	url        string
	projectKey string
}

type Project struct {
	Key string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

type CContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Content struct {
	Type    string     `json:"type"`
	Content []CContent `json:"content"`
}

type Description struct {
	Type    string    `json:"type"`
	Version int       `json:"version"`
	Content []Content `json:"content"`
}

type Fields struct {
	Project     Project     `json:"project"`
	Summary     string      `json:"summary"`
	IssueType   IssueType   `json:"issuetype"`
	Description Description `json:"description"`
}

type OutwardIssue struct {
	Key string `json:"key"`
}

type InwardIssue struct {
	Key string `json:"key"`
}

type Type struct {
	Name string `json:"name"`
}

// Create is payload struct for CreateIssue API
type Create struct {
	Fields Fields `json:"fields"`
}

// Link is payload struct for LinkIssues API
type Link struct {
	Type         Type         `json:"type"`
	InwardIssue  InwardIssue  `json:"inwardIssue"`
	OutwardIssue OutwardIssue `json:"outwardIssue"`
}

//////////////////////
//
// Public API
//
//////////////////////

func NewJira(url, projectKey string) *Jira {
	return &Jira{url: url, projectKey: projectKey}
}

// TODO:
func (j *Jira) GetTicket(id string) (string, error) {
	return "", errors.New("TO BE IMPLEMENTED")
}

func (j *Jira) CreateIssue(summary, desc string) (string, error) {
	// load env vars
	loadedVars, err := internal.LoadEnv("JIRA_API_KEY", "JIRA_EMAIL")
	if err != nil {
		return "", err
	}

	payload := Create{
		Fields{
			Project: Project{
				Key: j.projectKey,
			},
			Summary: summary,
			IssueType: IssueType{
				Name: "Task",
			},
			Description: Description{
				Type:    "doc",
				Version: 1,
				Content: []Content{
					{
						Type: "paragraph",
						Content: []CContent{
							{
								Type: "text",
								Text: desc,
							},
						},
					},
				},
			},
		},
	}

	// convert payload to bytes
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// prepare http request
	req, err := http.NewRequest("POST", j.url+"/rest/api/3/issue", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(loadedVars[1], loadedVars[0])
	req.Header.Add("Content-Type", "application/json")

	// do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// decode data
	var content map[string]any

	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return "", err
	}

	// if error status code, throw error
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to create issue:\n%s", content)
	}

	// get response's key and convert it to string
	key, ok := content["key"].(string)
	if !ok {
		return "", errors.New("error converting response to string")
	}

	return key, nil
}

func (j *Jira) LinkIssues(fromTicket, toTicket string) (string, error) {
	// load env vars
	loadedVars, err := internal.LoadEnv("JIRA_API_KEY", "JIRA_EMAIL")
	if err != nil {
		return "", err
	}

	payload := Link{
		Type: Type{
			Name: "Relates", // TODO: Change me ??
		},
		InwardIssue: InwardIssue{
			Key: fromTicket,
		},
		OutwardIssue: OutwardIssue{
			Key: toTicket,
		},
	}

	// convert payload to bytes
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// prepare http request
	req, err := http.NewRequest("POST", j.url+"/rest/api/3/issueLink", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(loadedVars[1], loadedVars[0])
	req.Header.Add("Content-Type", "application/json")

	// do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If success and no body expected, stop here
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		// Linking issues doesn’t return a key, so return success
		return "issues linked successfully", nil
	}

	// Otherwise try to decode response (error case)
	var content map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil && err != io.EOF {
		return "", err
	}

	// Handle Jira API errors
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to link issues: %v", content)
	}

	return "", errors.New("reached impossible dead end. Aborting")
}
