# read-it

`read-it` is a small command-line TUI ([bubbletea](https://github.com/charmbracelet/bubbletea)) tool written in Go.  
Its purpose is to help developers browse Cypress test files, pick a test, automatically generate a summary using **GitHub Copilot CLI** or **Cursor CLI**, review/adjust the summary, and then create a Jira “test” ticket linked to a bug or improvement.

### Disclaimer

**Copilot** or **Cursor** are supported. You must install and configure those tools on your machine yourself-`read-it` does not bundle them and only invokes `copilot` or `agent` on your `PATH`. If the CLI for your chosen agent is missing, the program exits with an error asking you to install it first.

![read-it-usage-gif](./assets/read-it.gif)

## Features

- Parses Cypress test files (JS/TS) and extracts all tests
- Presents a terminal UI allowing you to browse and choose a test to summarise
- Summarises tests via **Copilot** or **Cursor** CLI (`-a copilot` | `-a cursor`)
- Lets user review and manually refine the summary before committing
- After confirmation, creates a new “test” ticket in Jira and links it to an existing ticket (bug, improvement, etc.)

## Download & Run

Download the latest version from the [Releases](https://github.com/Daniel-Rodrigues-Nokia/read-it/releases) page.

## Requirements

- **Go** `1.24` or newer (for building from source; see `go.mod`)
- **For `-a copilot`**: GitHub Copilot CLI installed so `copilot --version` succeeds on your `PATH`
- **For `-a cursor`**: Cursor CLI installed so `agent --version` succeeds on your `PATH`
- **Jira**: API access and a `read-it.env` file (see below) with permission to create **Xray Test** issues, assign them, link with **Is a test for**, and transition to **Closed** (your workflow must expose a transition named **Closed**)

## Configuration

`read-it` loads Jira settings from **`read-it.env`** in your **user config directory** (`os.UserConfigDir()` - on Linux typically `~/.config/read-it.env`. If the file is missing on first load, a blank template is created and you must fill it in:

- `JIRA_API_KEY`
- `JIRA_USER`
- `JIRA_URL`
- `JIRA_PROJECT`

## Installation

```bash
# Clone the repo
git clone https://github.com/dffrs/read-it.git
cd read-it

# Build the binary
make

# Install into your personal bin (create the directory if it does not exist)
mkdir -p "$HOME/bin"
cp read-it "$HOME/bin/read-it"
# or: mv read-it "$HOME/bin/read-it"
```

Ensure `$HOME/bin` is on your `PATH` (for example in `~/.bashrc`: `export PATH="$HOME/bin:$PATH"`).

## Usage

Example:

```bash
read-it -f path/to/spec.cy.ts
```

```
Usage: read-it [options]

Options:
	-f <file>		Path to cypress test file
	-a <agent>		Agent to use: copilot | cursor. Default = copilot
	-debug <mode>	Enable debug output and stop execution early.
						Available modes:
							* selection:	After choosing a test, print the selected test details and exit.
							* textarea:		After reviewing AI summary, print the final output and exit.
```

## TODO:

- [x] ~~Improve JIRA step by:~~
  - [x] ~~Assigning JIRA test ticket to its creator~~
  - [x] ~~Automatically closing JIRA test ticket~~
  - [x] ~~Use 'Subject' as ticket's title~~
- [x] ~~Improve 'instructions.txt' to use proper formatting (bullet points)~~
- [x] ~~Support Cursor~~
- [x] ~~Gracefully handle errors~~
- [x] ~~Alert user for updates~~
- [ ] Multiple test selection
