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
	"strings"
)

//////////////////////
//
// Structs
//
//////////////////////

type Jira struct {
	*internal.Config
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
	Project     Project   `json:"project"`
	Summary     string    `json:"summary"`
	IssueType   IssueType `json:"issuetype"`
	Description string    `json:"description"`
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

type To struct {
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	ID          string `json:"id"`
	Self        string `json:"self"`
}

type Transition struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OpsBarSequence int    `json:"opsbarSequence"`
	To             To     `json:"to"`
}

type Transitions struct {
	Expand      string       `json:"type"`
	Transitions []Transition `json:"transitions"`
}

//////////////////////
//
// Public API
//
//////////////////////

func NewJira(config *internal.Config) *Jira {
	return &Jira{config}
}

// TODO:

func (j *Jira) GetTicket(id string) (string, error) {
	return "", errors.New("TO BE IMPLEMENTED")
}

func (j *Jira) CreateXrayTest(summary, desc string) (string, error) {
	payload := Create{
		Fields{
			Project: Project{
				Key: j.JiraProject,
			},
			Summary: summary,
			IssueType: IssueType{
				Name: "Xray Test",
			},
			Description: desc,
		},
	}

	// convert payload to bytes
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// prepare http request
	req, err := http.NewRequest("POST", j.JiraURL+"/rest/api/latest/issue", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+j.JiraAPIKey)
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
	payload := Link{
		Type: Type{
			Name: "Is a test for",
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
	req, err := http.NewRequest("POST", j.JiraURL+"/rest/api/latest/issueLink", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+j.JiraAPIKey)
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

func (j *Jira) GetAllTransitions(ticket, state string) (*Transition, error) {
	req, err := http.NewRequest("GET", j.JiraURL+"/rest/api/latest/issue/"+ticket+"/transitions", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+j.JiraAPIKey)
	req.Header.Add("Content-Type", "application/json")

	// do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode data
	var content Transitions

	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return nil, err
	}

	// if error status code, throw error
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get transitions\n%v", content)
	}

	var matchedTransition Transition

	for _, v := range content.Transitions {
		if strings.EqualFold(v.Name, state) {
			matchedTransition = v
			break
		}
	}

	if matchedTransition.Name == "" {
		return nil, fmt.Errorf("failed to find transition state %s", state)
	}

	return &matchedTransition, nil
}

func (j *Jira) TransitionTo(ticket, state string) error {
	matchedTransition, err := j.GetAllTransitions(ticket, state)
	if err != nil {
		return err
	}

	transition := map[string]any{"transition": map[string]string{
		"id": matchedTransition.ID, // id for 'Closed'
	}}

	data, err := json.Marshal(transition)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", j.JiraURL+"/rest/api/latest/issue/"+ticket+"/transitions", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+j.JiraAPIKey)
	req.Header.Add("Content-Type", "application/json")

	// do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If success and no body expected, stop here
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		// Transitioning issues doesn’t return a key, so return success
		return nil
	}

	// Otherwise try to decode response (error case)
	var content map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil && err != io.EOF {
		return err
	}

	// Handle Jira API errors
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to link issues: %v", content)
	}

	return errors.New("reached impossible dead end. Aborting")
}
