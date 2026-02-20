# read-it (wip)

`read-it` is a small, self-hosted command-line TUI ([bubbletea](https://github.com/charmbracelet/bubbletea)) tool written in Go.  
Its purpose is to help developers browse Cypress test files, pick a test, automatically generate a summary via an LLM (using your API key), review/adjust the summary, and then create a Jira “test” ticket linked to a bug or improvement.

![read-it-usage-gif](./assets/read-it.gif)

## Features

- Parses Cypress test files (JS/TS) and extracts all tests
- Presents a terminal UI allowing you to browse and choose a test to summarise
- Integrates with an LLM (via API key from `.env`)
- Lets user review and manually refine the summary before committing
- After confirmation, creates a new “test” ticket in Jira and links it to an existing ticket (bug, improvement, etc.)

## Requirements

- Access to a compatible LLM API (e.g. via API key in environment)
- Credentials / permissions to create tickets in your Jira instance

## Installation

```bash
# Clone the repo
git clone https://github.com/dffrs/read-it.git
cd read-it

# Build the binary
make

# Usage
Usage: read-it [options]

Options:
	-f <file>	    Path to cypress test file
	-debug <mode>	Enable debug output and stop execution early.
			          Available modes:
				          selection:  After choosing a test, print the selected test details and exit.
				          textarea:   After reviewing AI summary, print the final output and exit.
```

## TODO:

- [x] Improve JIRA step by:
  - [x] ~~Assigning JIRA test ticket to its creator~~
  - [x] ~~Automatically closing JIRA test ticket~~
  - [x] ~~Use 'Subject' as ticket's title~~
- [x] ~~Improve 'instructions.txt' to use proper formatting (bullet points)~~
- [ ] Gracefully handle errors
