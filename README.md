# UI-Elf : The UI Element finder CLI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Test & Lint](https://github.com/amasotti/ui-elf/actions/workflows/test-lint.yml/badge.svg)](https://github.com/amasotti/ui-elf/actions/workflows/test-lint.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/amasotti/ui-elf)


A command-line tool that scans Vue.js and React codebases to locate specific component types (forms, buttons, dialogs, and custom components). 
The tool helps development teams audit their frontend applications by identifying where components are used and providing usage statistics.

## What's inside

- **Multi-framework support**: Scans both Vue.js (.vue) and React (.jsx, .tsx) files
- **Design library recognition**: Identifies components from popular libraries (Quasar, Material UI)
- **Flexible filtering**: Exclude test files and node_modules, or specify directories to scan
- **Fast scanning**: Concurrent file processing for efficient codebase analysis

## Installation

### Prerequisites

- Go 1.25.0 or higher

### Building from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd ui-elf
```

2. Build the binary:
```bash
go build -o ui-elf cmd/ui-elf/main.go
```

3. (Optional) Install globally:
```bash
go install cmd/ui-elf/main.go
```

Or move the binary to your PATH:
```bash
sudo mv ui-elf /usr/local/bin/
```

## Usage

### Basic Syntax

```bash
ui-elf --component-type <type> --directory <path> [options]
```

### Command-Line Flags

| Flag | Short | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--component-type` | `-t` | Component type to search for: `form`, `button`, `dialog`, or `custom` | Yes | - |
| `--directory` | `-d` | Directory to scan | No | `.` (current directory) |
| `--filter` | `-f` | Comma-separated list of directories to include | No | All directories |
| `--output` | `-o` | Output format: `terminal`, `json`, or `both` | No | `terminal` |


## Supported Components

### Forms
- Native HTML: `<form>`
- Quasar: `<q-form>`
- Material UI: `<Form>`, `<MuiForm>`, `<v-form>`

### Buttons
- Native HTML: `<button>`
- Quasar: `<q-btn>`
- Material UI: `<Button>`, `<MuiButton>`, `<v-btn>`

### Dialogs
- Native HTML: `<dialog>`
- Quasar: `<q-dialog>`
- Material UI: `<Dialog>`, `<MuiDialog>`, `<v-dialog>`

### Custom Components
When using `--component-type custom`, the tool will identify all custom component usage in your codebase.


## File Filtering

The tool automatically excludes:
- `node_modules` directory
- Test files (files/directories containing: `test`, `tests`, `__tests__`, `.test.`, `.spec.`)

Use the `--filter` flag to scan only specific directories:
```bash
ui-elf -t form -d . -f src/components,src/views
```

## License

MIT License, see [LICENSE](LICENSE) for details.
