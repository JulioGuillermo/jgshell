# jgshell

`jgshell` is a Go-based TUI shell wrapper that runs your real shell inside a PTY and layers a structured interface on top of it. Instead of a plain prompt, it gives you a command timeline, syntax-highlighted input, completion and history pickers, persistent command history, and a live status bar with shell, directory, OS, and Git information.

## Features

### Terminal UI

- Full-screen terminal UI built with Bubble Tea and Lip Gloss.
- Command history rendered as scrollable cards instead of an unstructured terminal log.
- Each command card shows:
  - the original command with syntax highlighting
  - cleaned command output
  - exit status
  - elapsed execution time
  - user, shell, and working directory context
- Running commands display a live spinner.
- Full-screen terminal programs are passed through directly while they are active.
- Mouse-aware viewport support for scrolling command output.
- Automatic layout resizing when the terminal window changes size.

### Input Experience

- Syntax-highlighted command input powered by Tree-sitter using the Bash grammar.
- Multi-line input support.
- Inline history-based suggestions while you type.
- Suggestion acceptance from the keyboard.
- Autocomplete insertion at the current cursor position.
- Separate autocomplete and history menus built on Bubble Tea lists.

### Command Execution

- Runs your configured shell inside a PTY using `creack/pty`.
- Sends normal keystrokes and terminal control keys to the running process.
- Tracks command lifecycle using shell wrapper markers to detect start, finish, current user, and current directory.
- Captures exit codes and command duration for every executed command.
- Supports interrupting or re-wrapping the shell session from the UI.
- Cleans carriage returns, backspaces, tabs, and alternate-screen noise from stored output before rendering.

### History and Recall

- Persistent command history stored in `~/.jgshell/history`.
- Duplicate commands are deduplicated before saving.
- History is capped to the most recent 5000 entries.
- Prefix filtering for history search.
- Latest matching command used as the source of inline suggestions.

### Status Bar

- Status bar shows:
  - current user
  - current directory, with home-directory compaction
  - detected shell
  - detected OS
  - Git repository state when available
- Git parsing includes:
  - current local branch
  - remote tracking branch
  - ahead / behind counts
  - untracked file count
  - modified file count
  - staged file count
  - deleted file count
  - conflict count

### Shell Support

- Defaults to launching `bash` when no shell command is provided.
- Current support is focused on POSIX-style shells.
- Confirmed supported shells today:
  - `bash`
  - `zsh`
  - `sh`
  - `ksh`
  - other `sh`-compatible shells when their prompt and execution behavior remain compatible with the wrapper approach
- Detects the active shell from inside the session.
- Command wrapping, status loading, and autocomplete are currently documented around the POSIX shell path.
- `powershell`, `fish`, and possibly `nushell` are planned targets, but they are not supported yet and should be treated as work in progress.

## Architecture

The project is organized into separate layers so shell integration, UI, and feature logic stay decoupled.

- `app/`: Bubble Tea application and UI components.
- `controller/`: Orchestration layer connecting UI, executors, status loading, autocomplete, and persistence.
- `executor/`: PTY-backed command execution, fast command probes, and history tracking.
- `shell/`: Shell process startup and PTY management.
- `status/`: Status collection and Git parsing.
- `autocomplete/`: Shell-aware completion execution.
- `syntax/`: Tree-sitter highlighter.
- `wrapper/`: Command wrapping and output marker parsing.
- `persistence/`: Command history storage.

## Prerequisites

- Go `1.25.7` or newer
- A terminal with UTF-8 and true-color support
- A supported shell installed locally

## Installation

```bash
git clone https://github.com/julioguillermo/jgshell.git
cd jgshell
go mod download
go build -o jgshell main.go
```

## Usage

Run with the default shell:

```bash
./jgshell
```

Run with a specific shell:

```bash
./jgshell zsh
```

Run with a shell command and arguments:

```bash
./jgshell "bash --noprofile --norc"
```

## Keybindings

### Global

- `Ctrl+E` or `Alt+Shift+E`: quit the application
- `Ctrl+Shift+W` or `Alt+Shift+W`: re-wrap the current shell session

### While Idle

- `Enter`: execute the current input
- `Shift+Enter` or `Alt+Enter`: insert a newline into the input
- `Tab`: open autocomplete for the current line
- `Ctrl+R`: open history search
- `Ctrl+Space`: accept the inline suggestion
- `Right`: accept the inline suggestion when the cursor is at the end of the input

### While a Command Is Running

- Typed characters are forwarded to the underlying shell or running program
- Special keys forwarded include `Enter`, `Tab`, `Esc`, arrow keys, `Backspace`, `Space`, and `Ctrl+C`

## Development

Common tasks:

```bash
go fmt ./...
go vet ./...
go test ./...
go build -o jgshell main.go
```

Core entry points:

- [`main.go`](/home/jg/Documents/mispro/jgshell/main.go)
- [`app/infrastructure/app/app.go`](/home/jg/Documents/mispro/jgshell/app/infrastructure/app/app.go)
- [`controller/infrastructure/shell_controller.go`](/home/jg/Documents/mispro/jgshell/controller/infrastructure/shell_controller.go)
- [`executor/application/executor.go`](/home/jg/Documents/mispro/jgshell/executor/application/executor.go)
- [`status/infrastructure/status_loader.go`](/home/jg/Documents/mispro/jgshell/status/infrastructure/status_loader.go)
- [`autocomplete/infrastructure/autocomplete.go`](/home/jg/Documents/mispro/jgshell/autocomplete/infrastructure/autocomplete.go)

Additional project guidance lives in [`AGENTS.md`](/home/jg/Documents/mispro/jgshell/AGENTS.md).

## License

This project is fully free and open source under the MIT License. See [`LICENSE`](/home/jg/Documents/mispro/jgshell/LICENSE).
