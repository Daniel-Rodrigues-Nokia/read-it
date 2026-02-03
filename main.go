package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"

	in "read-it/internal"
	qu "read-it/internal/components/queue"
	sl "read-it/internal/components/selector"
	sp "read-it/internal/components/spinner"
	ta "read-it/internal/components/textarea"
	ti "read-it/internal/components/textinput"
	ji "read-it/internal/jira"
	rd "read-it/internal/reader"
)

//go:embed instructions.txt
var instructions string

const (
	debugTestSelection = "selection"
	debugTextareaOuput = "textarea"
)

func main() {
	filePath := flag.String("f", "", "path to file")
	debug := flag.String("debug", "", "debug")

	flag.Parse()

	// no file path ? throw error
	if *filePath == "" {
		usageMsg := in.GetUsageMsg()
		in.ThrowError(usageMsg)
	}

	config, err := in.LoadEnv()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// scan file and get tests
	tests, err := rd.ScanTests(*filePath, rd.Cypress{})
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

	// debug test selection stops execution here
	if *debug == debugTestSelection {
		fmt.Printf("Output:\nindex: %d\ntest:\n%s", optionsChosen, tests[optionsChosen].PrintItem())
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
	// resp, err := ai.AISummary(instructions, []string{tests[optionsChosen].PrintItem()}, ai.Copilot{})
	// if err != nil {
	// 	in.ThrowError(err.Error())
	// }

	resp := new([]string)
	*resp = append(*resp, "test")

	spinner.Stop()
	in.ClearStdOut()

	// -------------------------------------------------- Phase 3: Review summary ---------------------------------------------------
	// Now, output that response to a textarea, so that we can review it
	// and make changes if needed

	summary := *resp
	testValidated, err := ta.NewTextarea(summary[0]).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	// debug textarea stops execution here
	if *debug == debugTextareaOuput {
		fmt.Printf("Output:\n%s\n", testValidated)
		os.Exit(0)
	}

	// ------------------------------------------- Phase 4: Create & Link JIRA Tickets ----------------------------------------------
	// After phase 3, ask for the main ticket to link this test summary to
	srcTicket, err := ti.NewInput("ID of the bug/improvement/task...", in.CancelCtrl).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	j := ji.NewJira(config)

	// now, let's create a queue. This queue will have 2 tasks:
	// - create a 'test' JIRA ticket
	// - link it to the main ticket (got from phase 3)
	firstTask := qu.NewTask("Creating Xray Test...", func(m qu.Model) (any, error) {
		testTitle := fmt.Sprintf("Test for: %s", srcTicket)
		return j.CreateXrayTest(testTitle, testValidated)
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

	thirdTask := qu.NewTask("Closing Xray ticket...", func(m qu.Model) (any, error) {
		firstTask, err := m.GetResultFromTask(0)
		if err != nil {
			return nil, err
		}

		xrayTicket, ok := firstTask.Result.(string)
		if !ok {
			return nil, errors.New("cannot convert ticketID to string")
		}

		err = j.TransitionTo(xrayTicket, "Closed")

		return nil, err
	})

	// 'start' queue
	_, err = qu.NewQueue(firstTask, secondTask, thirdTask).Start()
	if err != nil {
		in.ThrowError(err.Error())
	}

	os.Exit(0)
}
