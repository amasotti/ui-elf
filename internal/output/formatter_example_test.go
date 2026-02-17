package output_test

import (
	"fmt"

	"ui-elf/internal/output"
	"ui-elf/internal/types"
)

// Example demonstrates terminal output formatting
func ExampleOutputFormatter_FormatTerminal() {
	formatter := output.NewOutputFormatter()

	result := &types.ScanResult{
		Matches: []types.ComponentMatch{
			{
				FilePath:      "src/components/UserForm.vue",
				Line:          12,
				ComponentName: "q-form",
				ComponentType: "form",
			},
			{
				FilePath:      "src/pages/Login.vue",
				Line:          25,
				ComponentName: "form",
				ComponentType: "form",
			},
		},
		TotalCount:    2,
		ScanTimeMs:    150,
		ComponentType: "form",
		ScannedFiles:  50,
	}

	output := formatter.FormatTerminal(result)
	fmt.Print(output)

	// Output:
	//
	// Component Finder Results - form
	// ==================================================
	//
	// Found components in:
	//
	//   src/components/UserForm.vue (line 12): q-form
	//   src/pages/Login.vue (line 25): form
	//
	// --------------------------------------------------
	// Total components found: 2
	// Files scanned: 50
	// Scan time: 150ms
}

// Example demonstrates JSON output formatting
func ExampleOutputFormatter_FormatJSON() {
	formatter := output.NewOutputFormatter()

	result := &types.ScanResult{
		Matches: []types.ComponentMatch{
			{
				FilePath:      "src/App.tsx",
				Line:          15,
				ComponentName: "Button",
				ComponentType: "button",
			},
		},
		TotalCount:    1,
		ScanTimeMs:    200,
		ComponentType: "button",
		ScannedFiles:  20,
	}

	jsonStr, _ := formatter.FormatJSON(result)
	fmt.Println(jsonStr)

	// Output:
	// {
	//   "matches": [
	//     {
	//       "filePath": "src/App.tsx",
	//       "line": 15,
	//       "componentName": "Button",
	//       "componentType": "button"
	//     }
	//   ],
	//   "totalCount": 1,
	//   "scanTimeMs": 200,
	//   "componentType": "button",
	//   "scannedFiles": 20
	// }
}
