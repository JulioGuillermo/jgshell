# Developer Guide: jgshell

## Project Overview

`jgshell` is a modern TUI shell wrapper built in Go. It aims to provide a rich terminal experience with syntax highlighting, interactive components, and shell-agnostic features (supporting bash, zsh, etc.).

## Technical Stack

- **Go 1.25.7**
- **TUI:** Charm Bubble Tea v2, Lipgloss v2
- **Syntax:** Tree-sitter (Go bindings)
- **Low-level:** PTY interaction via `creack/pty`

## Architecture

The project is organized into layers to separate concerns:

- **App Layer (`app/`):** Manages the Bubble Tea lifecycle and UI components.
- **Controller Layer (`controller/`):** The orchestrator. It connects the shell to the UI and coordinates executors.
- **Executor Layer (`executor/`):** Handles the actual interaction with the PTY, history management, and command output reading.
- **Domain vs Infrastructure:**
  - `domain/` contains interfaces and business logic.
  - `infrastructure/` contains concrete implementations (e.g., how to talk to a specific shell or render a specific component).

## Core Components

- **ShellController:** The central hub. It initializes the shell, executors, and features (status, autocomplete).
- **CmdViewPort:** A component for rendering command output with syntax highlighting.
- **Input:** A custom text input component that handles autocompletion and command submission.
- **StatusLoader:** Periodically gathers system information (Git status, current directory) to update the UI.

## Development Workflow

1. **Setup:** `go mod download`
2. **Build:** `go build -o jgshell main.go`
3. **Run:** `./jgshell`
4. **Test:** `go test ./...`
5. **Lint:** `go vet ./...` and `go fmt ./...`

## Key Packages

- `app/infrastructure/components/`: Sub-components like `cmdcard`, `input`, `menu`, `statusbar`.
- `syntax/infrastruct/`: Tree-sitter highlighter implementation.
- `wrapper/`: Logic for wrapping commands to detect their exit status and capture output precisely.

## Adding Features

- **UI:** Implement the `tea.Model` interface if it's a major component, or integrate into the `App` struct.
- **Executors:** If adding a new way to run commands, look at `executor/domain/executor.go`.
- **Status:** New status indicators should be added to `status/domain/status.go` and implemented in `status/infrastructure/`.

## Contributing

- Follow Go's official style guidelines.
- Ensure all exported symbols are documented.
- Maintain the separation between domain and infrastructure.
