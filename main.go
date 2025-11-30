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

	// right now, it supports Cypress and Java (untested) test cases
	tests, err := rd.ScanTests(path, cy.Cypress{})
	if err != nil {
		in.ThrowError(err.Error())
	}

	// ------------------------------------------------------- Phase 1: Get tests ------------------------------------------------------
	// Render a list of all the test detected
	// return the index of the chosen option
	optionsChosen, err := sl.NewSelector(tests, "Choose tests to summarize").Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// ExitOption ? Then abort this whole thing
	if optionsChosen == sl.ExitOption {
		os.Exit(0)
	}

	// ----------------------------------------------- Phase 2: Get AI to summarize them -----------------------------------------------
	// After that, let's get AI to summarize it
	// while we wait, we get a nice spinner animation :)
	spinner, err := sp.NewSpinner("Generating summary...", in.CancelCtrl).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// ('getting the summary' only starts here)
	resp, cancel, err := ai.SummarizeTests("instructions.txt", []string{tests[optionsChosen].PrintItem()})
	if err != nil {
		in.ThrowError(err.Error())
	}
	defer cancel()
	defer resp.Body.Close()

	spinner.Stop()

	// read AI's response
	content, err := ai.ReadSummary(resp)
	if err != nil {
		in.ThrowError(err.Error())
	}

	// do standard checks
	if len(content.Choices) < 1 {
		in.ThrowError(errors.New("got no response back").Error())
	}

	in.ClearStdOut()

	// -------------------------------------------------- Phase 3: Review summary ---------------------------------------------------
	// Now, output that response to a textarea, so that we can review it
	// and make changes if needed
	testValidated, err := ta.NewTextarea(content.Choices[0].Message.Content).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// ------------------------------------------- Phase 4: Create & Link JIRA Tickets ----------------------------------------------
	// After phase 3, ask for the main ticket to link this test summary to
	srcTicket, err := ti.NewInput("...", in.CancelCtrl).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// load JIRA variables
	loadedVar, err := in.LoadEnv("JIRA_URL", "JIRA_PROJECT")
	if err != nil {
		in.ThrowError(errors.New("jira url was not detected in .env file").Error())
	}

	j := ji.NewJira(loadedVar[0], loadedVar[1])

	// now, let's create a queue. This queue will have 2 tasks:
	// - create a 'test' JIRA ticket
	// - link it to the main ticket (got from phase 3)
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

	// 'start' queue
	_, err = qu.NewQueue(firstTask, secondTask).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	os.Exit(0)
}
