# jgshell

`jgshell` is a modern, high-performance terminal shell wrapper written in Go. It transforms your standard shell (bash, zsh, etc.) into a rich Text User Interface (TUI) experience featuring real-time syntax highlighting, interactive autocompletion, and integrated system status indicators.

## 🚀 Features

- **TUI-Powered Experience:** Built with the Charm [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for a smooth, interactive interface.
- **Smart Syntax Highlighting:** Leverages [Tree-sitter](https://tree-sitter.io/) for accurate, language-aware command highlighting as you type.
- **Interactive Autocompletion:** A dedicated menu for command and path suggestions.
- **Contextual Status Bar:** Real-time visibility into Git status, current working directory, and system environment.
- **Shell Agnostic:** Wraps your existing shell (bash by default) while maintaining your aliases and configurations.
- **High Performance:** Utilizes Go's concurrency and low-level PTY handling for near-zero latency.

## 🛠️ Architecture

`jgshell` follows a clean, hexagonal-inspired architecture:

- **App Layer:** Manages the TUI lifecycle and component rendering.
- **Controller Layer:** Orchestrates communication between the UI and the underlying shell.
- **Executor Layer:** Handles low-level PTY interaction, command wrapping, and output capturing.
- **Infrastructure Implementation:** Decoupled implementations for different shell environments and UI components.

## 📋 Prerequisites

- **Go:** 1.25.7 or higher.
- **Terminal:** A modern terminal with Unicode and True Color support (e.g., Alacritty, iTerm2, Kitty).

## 📥 Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/julioguillermo/jgshell.git
   cd jgshell
   ```

2. Download dependencies:

   ```bash
   go mod download
   ```

3. Build the binary:
   ```bash
   go build -o jgshell main.go
   ```

## 🎮 Usage

Simply run the compiled binary:

```bash
./jgshell
```

By default, it will start `bash`. You can specify a different shell as an argument:

```bash
./jgshell zsh
```

## 🧪 Development

### Key Commands

- **Format code:** `go fmt ./...`
- **Run linter:** `go vet ./...`
- **Run tests:** `go test ./...`
- **Build:** `go build -o jgshell main.go`

For detailed developer instructions, please refer to [AGENTS.md](./AGENTS.md) and [GEMINI.md](./GEMINI.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
