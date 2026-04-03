# AGENTS.md - Developer Guidelines for jgshell

## Project Overview

jgshell is a Go-based terminal shell application using the Charm Bubble Tea framework. It provides a TUI shell interface with syntax highlighting via tree-sitter.

## Build, Lint & Test Commands

```bash
# Build the application
go build -v ./...

# Run all tests (if any exist)
go test ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Run with coverage
go test -v -cover ./...

# Format code
go fmt ./...

# Run go vet
go vet ./...

# Run all checks (fmt + vet + staticcheck)
go fmt ./... && go vet ./...

# Install dependencies
go mod download

# Update dependencies
go get -u ./... && go mod tidy
```

## Code Style Guidelines

### General Principles

- **Keep it simple**: Prefer readable, straightforward code over clever abstractions
- **Small functions**: Break down complex logic into focused, single-responsibility functions
- **Explicit over implicit**: Make intentions clear, avoid magic behavior

### Imports

```go
import (
    "fmt"
    "os"

    tea "charm.land/bubbletea/v2"  // external packages with alias
    "github.com/julioguillermo/jgshell/app/infrastructure/app"  // local packages
)
```

- Group imports: standard library first, then external packages, then local packages
- Use blank line separation between groups
- Use aliases for packages with conflicting names (e.g., `tea`, `lipgloss`)
- Avoid relative imports like `./package`

### Formatting

- Use `go fmt` or goimports for automatic formatting
- Maximum line length: 100 characters (soft guideline)
- Use tabs for indentation, not spaces
- Use blank lines generously to separate logical sections within functions

### Types

```go
// Prefer concrete types over interfaces unless abstraction is needed
type State struct {
    shell shell.Shell
    cmds  []domain.Cmd
}

// Use pointers for large structs or when mutation is needed
func NewState(s shell.Shell) *State {
    return &State{shell: s}
}

// Return interfaces only when polymorphism is required
```

### Naming Conventions

- **Files**: `snake_case.go` (e.g., `state_cmd.go`, `app_update.go`)
- **Types**: `PascalCase` (e.g., `State`, `App`, `CmdCard`)
- **Functions/Variables**: `camelCase` (e.g., `NewState`, `GetCommands`)
- **Constants**: `PascalCase` or `camelCase` for private constants
- **Interfaces**: `PascalCase`, often with `er` suffix (e.g., `Reader`, `Runner`)
- **Private methods/fields**: start with lowercase (e.g., `update`, `shell`)
- **Acronyms**: use all caps for 2-letter acronyms (e.g., `URL`, `ID`), camelCase for 3+ (e.g., `url`, `html`)

### Error Handling

```go
// Always handle errors explicitly
shell, err := shellinfrastructure.NewShellConnector("bash")
if err != nil {
    fmt.Printf("Fail to start shell: %v", err)
    os.Exit(1)
}

// For functions returning (value, error), check error before using value
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}
```

- Never ignore returned errors with `_`
- Use `fmt.Errorf` with `%w` for wrapping errors
- For fatal errors in main/startup, `os.Exit(1)` is acceptable
- For recoverable errors, return error for caller to handle

### Packages & Directory Structure

```
app/infrastructure/     # UI components and infrastructure
shell/                  # Shell connector
state/application/      # Application state
state/domain/           # Domain models and business logic
syntax/                 # Syntax highlighting
tools/                  # Utility programs
```

- Each package should have a focused purpose
- Use descriptive package names
- Avoid circular dependencies between packages

### Concurrency

- Use goroutines for concurrent operations
- Always provide a way to stop goroutines (channels, context)
- Lock mutexes for shared state access
- Prefer structured concurrency patterns over sleep-based timing

### TUI-Specific Guidelines (Bubble Tea)

```go
// Models implement tea.Model interface
type App struct {
    state  *application.State
    highlighter *syntaxinfrastructure.TSHighlighter
}

// Messages trigger updates
type TickMsg time.Time

// Update function signature
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case TickMsg:
        // handle tick
    }
    return a, nil
}

// View renders the UI
func (a *App) View() string {
    return someLipglossStyle.Render("content")
}
```

- Keep update functions fast; do heavy work in commands
- Use lipgloss for styling (this project uses `charm.land/lipgloss/v2`)
- Commands (tea.Cmd) handle side effects and async work

### Testing

- Create test files with `_test.go` suffix in the same package
- Use table-driven tests for multiple test cases
- Name tests descriptively: `TestFunctionName_ExpectedBehavior`
- Test behavior, not implementation details

### Git Conventions

- Commit messages: concise, imperative mood ("Add feature" not "Added feature")
- Keep commits atomic and focused
- Run `go fmt` and `go vet` before committing

## Dependencies

- **Go version**: 1.25.7
- **Key dependencies**:
  - `charm.land/bubbletea/v2` - TUI framework
  - `charm.land/lipgloss/v2` - Terminal styling
  - `github.com/tree-sitter/go-tree-sitter` - Syntax highlighting
  - `github.com/creack/pty` - PTY handling

## Notes for Agents

- This is a TUI application - verify UI renders correctly
- The binary builds to `./jgshell` (root directory)
- No existing test suite - consider adding tests for new functionality
- Run `go mod tidy` after adding dependencies
