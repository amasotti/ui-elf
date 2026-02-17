package registry

import "testing"

func TestNewComponentMappingRegistry(t *testing.T) {
	registry := NewComponentMappingRegistry()

	if registry == nil {
		t.Fatal("Expected registry to be created, got nil")
	}

	if len(registry.mappings) == 0 {
		t.Fatal("Expected registry to have mappings, got empty map")
	}
}

func TestGetMapping(t *testing.T) {
	registry := NewComponentMappingRegistry()

	tests := []struct {
		name          string
		componentType string
		shouldExist   bool
	}{
		{"form mapping exists", "form", true},
		{"button mapping exists", "button", true},
		{"dialog mapping exists", "dialog", true},
		{"unknown mapping", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := registry.GetMapping(tt.componentType)
			if exists != tt.shouldExist {
				t.Errorf("GetMapping(%q) exists = %v, want %v", tt.componentType, exists, tt.shouldExist)
			}
		})
	}
}

func TestMatchesComponentType_Forms(t *testing.T) {
	registry := NewComponentMappingRegistry()

	tests := []struct {
		name          string
		componentName string
		shouldMatch   bool
	}{
		{"native form", "form", true},
		{"quasar q-form", "q-form", true},
		{"quasar QForm", "QForm", true},
		{"material v-form", "v-form", true},
		{"material VForm", "VForm", true},
		{"material Form", "Form", true},
		{"material MuiForm", "MuiForm", true},
		{"case insensitive", "FORM", true},
		{"case insensitive quasar", "Q-FORM", true},
		{"non-form component", "button", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := registry.MatchesComponentType(tt.componentName, "form")
			if matches != tt.shouldMatch {
				t.Errorf("MatchesComponentType(%q, %q) = %v, want %v",
					tt.componentName, "form", matches, tt.shouldMatch)
			}
		})
	}
}

func TestMatchesComponentType_Buttons(t *testing.T) {
	registry := NewComponentMappingRegistry()

	tests := []struct {
		name          string
		componentName string
		shouldMatch   bool
	}{
		{"native button", "button", true},
		{"quasar q-btn", "q-btn", true},
		{"quasar QBtn", "QBtn", true},
		{"material v-btn", "v-btn", true},
		{"material VBtn", "VBtn", true},
		{"material Button", "Button", true},
		{"material MuiButton", "MuiButton", true},
		{"case insensitive", "BUTTON", true},
		{"non-button component", "form", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := registry.MatchesComponentType(tt.componentName, "button")
			if matches != tt.shouldMatch {
				t.Errorf("MatchesComponentType(%q, %q) = %v, want %v",
					tt.componentName, "button", matches, tt.shouldMatch)
			}
		})
	}
}

func TestMatchesComponentType_Dialogs(t *testing.T) {
	registry := NewComponentMappingRegistry()

	tests := []struct {
		name          string
		componentName string
		shouldMatch   bool
	}{
		{"native dialog", "dialog", true},
		{"quasar q-dialog", "q-dialog", true},
		{"quasar QDialog", "QDialog", true},
		{"material v-dialog", "v-dialog", true},
		{"material VDialog", "VDialog", true},
		{"material Dialog", "Dialog", true},
		{"material MuiDialog", "MuiDialog", true},
		{"case insensitive", "DIALOG", true},
		{"non-dialog component", "button", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := registry.MatchesComponentType(tt.componentName, "dialog")
			if matches != tt.shouldMatch {
				t.Errorf("MatchesComponentType(%q, %q) = %v, want %v",
					tt.componentName, "dialog", matches, tt.shouldMatch)
			}
		})
	}
}

func TestMatchesComponentType_CustomComponent(t *testing.T) {
	registry := NewComponentMappingRegistry()

	// For custom component types (not in registry), should do exact name match
	tests := []struct {
		name          string
		componentName string
		componentType string
		shouldMatch   bool
	}{
		{"exact match", "MyCustomComponent", "MyCustomComponent", true},
		{"case insensitive match", "mycustomcomponent", "MyCustomComponent", true},
		{"no match", "OtherComponent", "MyCustomComponent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := registry.MatchesComponentType(tt.componentName, tt.componentType)
			if matches != tt.shouldMatch {
				t.Errorf("MatchesComponentType(%q, %q) = %v, want %v",
					tt.componentName, tt.componentType, matches, tt.shouldMatch)
			}
		})
	}
}
