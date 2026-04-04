# Project Mandates: jgshell

This file contains foundational mandates for the Gemini CLI when working on the `jgshell` project. These rules take precedence over general defaults.

## Architectural Integrity

- **Pattern:** Adhere to the hexagonal-inspired layered architecture (`domain` vs `infrastructure`).
- **TUI:** All UI changes must be made within `app/infrastructure/components` using `charm.land/bubbletea/v2` and `lipgloss/v2`.
- **Logic Separation:** Business logic belongs in `domain` packages. I/O and external interactions belong in `infrastructure`.

## Testing & Validation

- **New Features:** Every new feature MUST include unit tests in the same package (e.g., `cmd_wrapper_test.go`).
- **Regression:** Run `go test ./...` before considering any task complete.
- **TUI Validation:** Since automated TUI testing is limited, manually verify UI changes by building and running `./jgshell` if possible, or by carefully reviewing the `View()` and `Update()` functions.

## Tool Usage

- **Go Commands:** Use `go fmt ./...`, `go vet ./...`, and `go mod tidy` as standard validation steps.
- **Tree-sitter:** When working with `syntax/`, ensure compatibility with the `github.com/tree-sitter/go-tree-sitter` API.

## Coding Style

- **Errors:** Use `fmt.Errorf("context: %w", err)` for error wrapping.
- **Naming:** Follow the project's established conventions: `PascalCase` for types, `camelCase` for variables, `snake_case` for files.
- **Indentation:** Use **Tabs** for Go files (standard `go fmt`).

## Security

- **PTY Handling:** Be extremely careful with `github.com/creack/pty` to avoid resource leaks; ensure `Close()` is called where appropriate.
- **Secrets:** Never log or display environment variables captured from the shell.
