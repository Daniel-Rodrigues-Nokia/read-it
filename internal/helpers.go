// Package internal helper functions
package internal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
)

const (
	HeaderColor string = "#f72e7f"
	HelpColor   string = "#4f4f4f"
	FocusColor  string = "#5eeddc"
	BlurColor   string = "#666666"
	MiscColor   string = "#ffffff"
)

func ThrowError(msg string) {
	fmt.Printf("%s\n", fmt.Errorf("[read-it]: %s", msg))
	os.Exit(1)
}

func LoadEnv(vars ...string) ([]string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var loadedVars []string
	for _, v := range vars {
		loadedVars = append(loadedVars, os.Getenv(v))
	}

	return loadedVars, nil
}

func ClearStdOut() {
	clearer := make(map[string]func())

	// linux based
	linuxCleanUp := func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	// for linux based machines
	clearer["linux"] = linuxCleanUp
	clearer["darwin"] = linuxCleanUp

	// for windows machines
	clearer["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	// clear stdout
	// if unsupported runtime, fail silently
	c, ok := clearer[runtime.GOOS]
	if ok {
		c()
	}
}

func CancelCtrl() {
	ClearStdOut()

	fmt.Print("\nProgram terminated\n")

	os.Exit(0)
}

func GetUsageMsg() string {
	return `
Usage: read-it [options]

Options:
	-f <file>	Path to cypress test file
	-debug <mode>	Enable debug output and stop execution early.
			Available modes:
				* selection:	After choosing a test, print the selected test details and exit.
				* textarea:	After reviewing AI summary, print the final output and exit.
	`
}
