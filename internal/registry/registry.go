// Package registry maintains mappings between component types and library-specific implementations.
package registry

import "strings"

// ComponentMapping defines the mapping structure for a component type
type ComponentMapping struct {
	Type     string
	Patterns map[string][]string // library name -> component names
}

// ComponentMappingRegistry manages mappings between component types and actual component names
type ComponentMappingRegistry struct {
	mappings map[string]ComponentMapping
}

// NewComponentMappingRegistry creates a new registry with hardcoded mappings
func NewComponentMappingRegistry() *ComponentMappingRegistry {
	registry := &ComponentMappingRegistry{
		mappings: make(map[string]ComponentMapping),
	}

	// Form mappings
	registry.mappings["form"] = ComponentMapping{
		Type: "form",
		Patterns: map[string][]string{
			"native":   {"form"},
			"quasar":   {"q-form", "QForm"},
			"material": {"v-form", "VForm", "Form", "MuiForm"},
		},
	}

	// Button mappings
	registry.mappings["button"] = ComponentMapping{
		Type: "button",
		Patterns: map[string][]string{
			"native":   {"button"},
			"quasar":   {"q-btn", "QBtn"},
			"material": {"v-btn", "VBtn", "Button", "MuiButton"},
		},
	}

	// Dialog mappings
	registry.mappings["dialog"] = ComponentMapping{
		Type: "dialog",
		Patterns: map[string][]string{
			"native":   {"dialog"},
			"quasar":   {"q-dialog", "QDialog"},
			"material": {"v-dialog", "VDialog", "Dialog", "MuiDialog"},
		},
	}

	return registry
}

// GetMapping returns the component mapping for a given component type
func (r *ComponentMappingRegistry) GetMapping(componentType string) (ComponentMapping, bool) {
	mapping, exists := r.mappings[strings.ToLower(componentType)]
	return mapping, exists
}

// MatchesComponentType checks if a component name matches a given component type
func (r *ComponentMappingRegistry) MatchesComponentType(componentName string, componentType string) bool {
	mapping, exists := r.GetMapping(componentType)
	if !exists {
		// For custom component types, do exact name match
		return strings.EqualFold(componentName, componentType)
	}

	// Check all patterns for the component type
	for _, patterns := range mapping.Patterns {
		for _, pattern := range patterns {
			if strings.EqualFold(componentName, pattern) {
				return true
			}
		}
	}

	return false
}
