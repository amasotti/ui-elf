# UI-Elf : The UI Element finder CLI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


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

### Examples

#### Scan for forms in current directory
```bash
ui-elf --component-type form --directory .
```

#### Scan for buttons with JSON output
```bash
ui-elf --component-type button --directory ./src --output json
```

#### Scan specific directories only
```bash
ui-elf --component-type dialog --directory . --filter src/components,src/views
```

#### Export results to both terminal and JSON
```bash
ui-elf --component-type form --directory . --output both
```

#### Scan for custom components
```bash
ui-elf --component-type custom --directory ./src
```

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

## Output Formats

### Terminal Output

```
Scanning for components...
Found 15 form components in 450 files (scan time: 1.234s)

Results:
  src/components/UserForm.vue:12 - q-form
  src/components/LoginForm.vue:8 - form
  src/views/Settings.jsx:45 - Form
  ...

Total: 15 components found
```

### JSON Output

Results are saved to `ui-elf-results.json`:

```json
{
  "componentType": "form",
  "totalCount": 15,
  "scanTimeMs": 1234,
  "scannedFiles": 450,
  "matches": [
    {
      "filePath": "src/components/UserForm.vue",
      "line": 12,
      "componentName": "q-form",
      "componentType": "form"
    }
  ]
}
```

## File Filtering

The tool automatically excludes:
- `node_modules` directory
- Test files (files/directories containing: `test`, `tests`, `__tests__`, `.test.`, `.spec.`)

Use the `--filter` flag to scan only specific directories:
```bash
ui-elf -t form -d . -f src/components,src/views
```

## Error Handling

The tool provides clear error messages for common issues:

- **Directory not found**: Displays error and exits
- **Invalid component type**: Shows valid options (form, button, dialog, custom)
- **Invalid output format**: Shows valid options (terminal, json, both)
- **Parse errors**: Logs warnings and continues scanning other files
- **No files found**: Displays "No files to scan" message
- **No components found**: Displays "0 components found" message

## Development

### Project Structure

```
ui-elf/
├── cmd/
│   └── ui-elf/
│       └── main.go           # Entry point
├── internal/
│   ├── cli/                  # CLI controller
│   ├── discovery/            # File discovery service
│   ├── output/               # Output formatting
│   ├── registry/             # Component mapping registry
│   ├── scanner/              # Component parsers and scanner
│   └── types/                # Core data structures
├── sample-files/             # Test files
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ui-elf-linux cmd/ui-elf/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o ui-elf-macos cmd/ui-elf/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o ui-elf.exe cmd/ui-elf/main.go
```

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
