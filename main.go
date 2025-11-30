package main

import (
	"errors"
	"fmt"
	"os"

	in "read-it/internal"
	ai "read-it/internal/aiSum"
	qu "read-it/internal/components/queue"
	sl "read-it/internal/components/selector"
	sp "read-it/internal/components/spinner"
	ta "read-it/internal/components/textarea"
	ti "read-it/internal/components/textinput"
	ji "read-it/internal/jira"
	rd "read-it/internal/reader"
	cy "read-it/internal/reader/cypress"
)

func main() {
	// file path as arg
	args := os.Args

	if len(args) < 2 {
		in.ThrowError("File path needed")
	}
	path := os.Args[1]

	tests, err := rd.BuildTests(path, cy.Cypress{})
	// tests, err := r.BuildTests(path, ju.JUnit{})
	if err != nil {
		in.ThrowError(err.Error())
	}

	optionsChosen, err := sl.NewSelector(tests, "Choose tests to summarize").Choose()
	if err != nil {
		in.ThrowError(err.Error())
	}

	if optionsChosen == sl.ExitOption {
		os.Exit(0)
	}

	spinner, err := sp.BuildSpinner("Generating summary...", in.CancelCtrl).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	resp, cancel, err := ai.SummarizeTests([]string{tests[optionsChosen].PrintItem()})
	if err != nil {
		in.ThrowError(err.Error())
	}
	defer cancel()
	defer resp.Body.Close()

	spinner.Stop()

	content, err := ai.ReadSummary(resp)
	if err != nil {
		in.ThrowError(err.Error())
	}

	if len(content.Choices) < 1 {
		in.ThrowError(errors.New("got no response back").Error())
	}

	in.ClearStdOut()

	testValidated, err := ta.NewValidator(content.Choices[0].Message.Content).Validate()
	if err != nil {
		in.ThrowError(err.Error())
	}

	srcTicket, err := ti.NewInput("BP-xxxx", in.CancelCtrl).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	loadedVar, err := in.LoadEnv("JIRA_URL", "JIRA_PROJECT")
	if err != nil {
		in.ThrowError(errors.New("jira url was not detected in .env file").Error())
	}

	j := ji.NewJira(loadedVar[0], loadedVar[1])

	firstTask := qu.NewTask("Creating JIRA ticket...", func(m qu.Model) (any, error) {
		testTitle := fmt.Sprintf("Test for: %s", srcTicket)
		return j.CreateIssue(testTitle, testValidated)
	})

	secondTask := qu.NewTask("Linking Issues...", func(m qu.Model) (any, error) {
		firstTask, err := m.GetResultFromTask(0)
		if err != nil {
			return nil, err
		}

		destTicket, ok := firstTask.Result.(string)
		if !ok {
			return nil, errors.New("cannot convert ticketID to string")
		}

		_, err = j.LinkIssues(srcTicket, destTicket)

		return nil, err
	})

	_, err = qu.NewQueue(firstTask, secondTask).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	os.Exit(0)
}
