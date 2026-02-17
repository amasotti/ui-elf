package scanner

import (
	"testing"
)

func TestVueParser_SupportsFile(t *testing.T) {
	parser := NewVueParser()
	
	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"vue file", "component.vue", true},
		{"vue file with path", "src/components/MyComponent.vue", true},
		{"uppercase extension", "Component.VUE", true},
		{"jsx file", "component.jsx", false},
		{"tsx file", "component.tsx", false},
		{"js file", "component.js", false},
		{"no extension", "component", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.SupportsFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("SupportsFile(%q) = %v, want %v", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestVueParser_Parse_TemplateSection(t *testing.T) {
	parser := NewVueParser()
	
	tests := []struct {
		name          string
		content       string
		expectedCount int
		expectedNames []string
	}{
		{
			name: "single component",
			content: `<template>
  <q-form>
    <input />
  </q-form>
</template>`,
			expectedCount: 1,
			expectedNames: []string{"q-form"},
		},
		{
			name: "multiple components",
			content: `<template>
  <q-form>
    <q-btn label="Submit" />
    <q-dialog v-model="show">
      <div>Content</div>
    </q-dialog>
  </q-form>
</template>`,
			expectedCount: 3,
			expectedNames: []string{"q-form", "q-btn", "q-dialog"},
		},
		{
			name: "self-closing components",
			content: `<template>
  <MyComponent />
  <AnotherComponent/>
</template>`,
			expectedCount: 2,
			expectedNames: []string{"MyComponent", "AnotherComponent"},
		},
		{
			name: "kebab-case components",
			content: `<template>
  <my-custom-component />
  <another-component></another-component>
</template>`,
			expectedCount: 2,
			expectedNames: []string{"my-custom-component", "another-component"},
		},
		{
			name: "mixed case components",
			content: `<template>
  <QForm>
    <QBtn />
    <q-dialog />
  </QForm>
</template>`,
			expectedCount: 3,
			expectedNames: []string{"QForm", "QBtn", "q-dialog"},
		},
		{
			name: "ignore HTML tags",
			content: `<template>
  <div>
    <span>Text</span>
    <button>Click</button>
    <form>
      <input />
    </form>
  </div>
</template>`,
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name: "components with attributes",
			content: `<template>
  <q-form @submit="onSubmit" class="my-form">
    <q-btn label="Submit" color="primary" />
  </q-form>
</template>`,
			expectedCount: 2,
			expectedNames: []string{"q-form", "q-btn"},
		},
		{
			name: "empty template",
			content: `<template>
</template>`,
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name: "no template section",
			content: `<script>
export default {
  name: 'MyComponent'
}
</script>`,
			expectedCount: 0,
			expectedNames: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := parser.Parse(tt.content, "test.vue")
			
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			
			if len(matches) != tt.expectedCount {
				t.Errorf("Parse() returned %d matches, want %d", len(matches), tt.expectedCount)
			}
			
			// Check component names
			for i, expectedName := range tt.expectedNames {
				if i >= len(matches) {
					t.Errorf("Missing match for component %q", expectedName)
					continue
				}
				if matches[i].ComponentName != expectedName {
					t.Errorf("Match %d: got component name %q, want %q", 
						i, matches[i].ComponentName, expectedName)
				}
			}
		})
	}
}

func TestVueParser_Parse_ScriptSection_JSX(t *testing.T) {
	parser := NewVueParser()
	
	tests := []struct {
		name          string
		content       string
		expectedCount int
		expectedNames []string
	}{
		{
			name: "JSX in render function",
			content: `<template>
  <div>Template</div>
</template>
<script>
export default {
  render() {
    return <QForm>
      <QBtn />
    </QForm>
  }
}
</script>`,
			expectedCount: 2, // 2 from JSX (QForm, QBtn) - div is HTML tag
			expectedNames: []string{"QForm", "QBtn"},
		},
		{
			name: "JSX in setup script",
			content: `<script setup>
const MyComponent = () => {
  return <CustomComponent />
}
</script>`,
			expectedCount: 1,
			expectedNames: []string{"CustomComponent"},
		},
		{
			name: "multiple JSX components",
			content: `<script>
export default {
  render() {
    return (
      <Dialog>
        <Button onClick={handleClick} />
        <Form onSubmit={handleSubmit}>
          <Input />
        </Form>
      </Dialog>
    )
  }
}
</script>`,
			expectedCount: 4,
			expectedNames: []string{"Dialog", "Button", "Form", "Input"},
		},
		{
			name: "no JSX in script",
			content: `<script>
export default {
  data() {
    return {
      message: 'Hello'
    }
  }
}
</script>`,
			expectedCount: 0,
			expectedNames: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := parser.Parse(tt.content, "test.vue")
			
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			
			if len(matches) != tt.expectedCount {
				t.Errorf("Parse() returned %d matches, want %d", len(matches), tt.expectedCount)
			}
			
			// Check component names
			for i, expectedName := range tt.expectedNames {
				if i >= len(matches) {
					t.Errorf("Missing match for component %q", expectedName)
					continue
				}
				if matches[i].ComponentName != expectedName {
					t.Errorf("Match %d: got component name %q, want %q", 
						i, matches[i].ComponentName, expectedName)
				}
			}
		})
	}
}

func TestVueParser_Parse_LineNumbers(t *testing.T) {
	parser := NewVueParser()
	
	content := `<template>
  <div>
    <q-form>
      <q-btn />
    </q-form>
  </div>
</template>

<script>
export default {
  name: 'TestComponent'
}
</script>`
	
	matches, err := parser.Parse(content, "test.vue")
	
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	
	// We expect to find q-form and q-btn (div is HTML tag and should be ignored)
	if len(matches) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(matches))
	}
	
	// Check that line numbers are reasonable (not 0 or negative)
	for i, match := range matches {
		if match.Line <= 0 {
			t.Errorf("Match %d: line number %d should be positive", i, match.Line)
		}
		
		// q-form should be on line 3, q-btn on line 4
		// (accounting for template start line)
		if match.ComponentName == "q-form" && match.Line != 3 {
			t.Errorf("q-form should be on line 3, got line %d", match.Line)
		}
		if match.ComponentName == "q-btn" && match.Line != 4 {
			t.Errorf("q-btn should be on line 4, got line %d", match.Line)
		}
	}
}

func TestVueParser_Parse_FilePath(t *testing.T) {
	parser := NewVueParser()
	
	content := `<template>
  <q-form />
</template>`
	
	filePath := "src/components/MyForm.vue"
	matches, err := parser.Parse(content, filePath)
	
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	
	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}
	
	if matches[0].FilePath != filePath {
		t.Errorf("FilePath = %q, want %q", matches[0].FilePath, filePath)
	}
}

func TestVueParser_Parse_ComplexRealWorld(t *testing.T) {
	parser := NewVueParser()
	
	content := `<template>
  <q-page class="flex flex-center">
    <q-form @submit="onSubmit" class="q-gutter-md">
      <q-input
        v-model="name"
        label="Name"
        :rules="[val => !!val || 'Field is required']"
      />
      
      <q-select
        v-model="option"
        :options="options"
        label="Select Option"
      />
      
      <div class="row">
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn label="Cancel" @click="onCancel" flat />
      </div>
    </q-form>
    
    <q-dialog v-model="showDialog">
      <q-card>
        <q-card-section>
          <div class="text-h6">Confirmation</div>
        </q-card-section>
        
        <q-card-actions align="right">
          <q-btn flat label="OK" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
export default {
  name: 'MyForm',
  data() {
    return {
      name: '',
      option: null,
      showDialog: false,
      options: ['Option 1', 'Option 2']
    }
  },
  methods: {
    onSubmit() {
      this.showDialog = true
    },
    onCancel() {
      this.$router.back()
    }
  }
}
</script>`
	
	matches, err := parser.Parse(content, "test.vue")
	
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	
	// Expected components (excluding HTML tags like div):
	// q-page, q-form, q-input, q-select, q-btn (2x), q-dialog, q-card, 
	// q-card-section, q-card-actions, q-btn (in dialog)
	expectedComponents := []string{
		"q-page", "q-form", "q-input", "q-select", "q-btn", "q-btn",
		"q-dialog", "q-card", "q-card-section", "q-card-actions", "q-btn",
	}
	
	// Debug: print what we found
	if len(matches) != len(expectedComponents) {
		t.Logf("Found components:")
		for i, match := range matches {
			t.Logf("  %d: %s (line %d)", i, match.ComponentName, match.Line)
		}
	}
	
	if len(matches) != len(expectedComponents) {
		t.Errorf("Expected %d components, got %d", len(expectedComponents), len(matches))
	}
	
	// Verify all expected components are found
	foundComponents := make(map[string]int)
	for _, match := range matches {
		foundComponents[match.ComponentName]++
	}
	
	// Check q-btn appears 3 times
	if foundComponents["q-btn"] != 3 {
		t.Errorf("Expected 3 q-btn components, got %d", foundComponents["q-btn"])
	}
	
	// Check other components appear once
	singleComponents := []string{"q-page", "q-form", "q-input", "q-select", 
		"q-dialog", "q-card", "q-card-section", "q-card-actions"}
	for _, comp := range singleComponents {
		if foundComponents[comp] != 1 {
			t.Errorf("Expected 1 %s component, got %d", comp, foundComponents[comp])
		}
	}
}

func TestExtractTemplateSection(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedContent   string
		expectedStartLine int
	}{
		{
			name: "basic template",
			content: `<template>
  <div>Content</div>
</template>`,
			expectedContent:   "\n  <div>Content</div>\n",
			expectedStartLine: 1,
		},
		{
			name: "template with lang attribute",
			content: `<template lang="pug">
  div Content
</template>`,
			expectedContent:   "\n  div Content\n",
			expectedStartLine: 1,
		},
		{
			name: "template after script",
			content: `<script>
export default {}
</script>

<template>
  <div>Content</div>
</template>`,
			expectedContent:   "\n  <div>Content</div>\n",
			expectedStartLine: 5,
		},
		{
			name:              "no template",
			content:           `<script>export default {}</script>`,
			expectedContent:   "",
			expectedStartLine: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, startLine := extractTemplateSection(tt.content)
			
			if content != tt.expectedContent {
				t.Errorf("extractTemplateSection() content = %q, want %q", 
					content, tt.expectedContent)
			}
			
			if startLine != tt.expectedStartLine {
				t.Errorf("extractTemplateSection() startLine = %d, want %d", 
					startLine, tt.expectedStartLine)
			}
		})
	}
}

func TestExtractScriptSection(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedHasScript bool
		expectedStartLine int
	}{
		{
			name: "basic script",
			content: `<script>
export default {}
</script>`,
			expectedHasScript: true,
			expectedStartLine: 1,
		},
		{
			name: "script setup",
			content: `<script setup>
const msg = 'Hello'
</script>`,
			expectedHasScript: true,
			expectedStartLine: 1,
		},
		{
			name: "script with lang",
			content: `<script lang="ts">
export default {}
</script>`,
			expectedHasScript: true,
			expectedStartLine: 1,
		},
		{
			name: "script after template",
			content: `<template>
  <div>Content</div>
</template>

<script>
export default {}
</script>`,
			expectedHasScript: true,
			expectedStartLine: 5,
		},
		{
			name:              "no script",
			content:           `<template><div>Content</div></template>`,
			expectedHasScript: false,
			expectedStartLine: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, startLine := extractScriptSection(tt.content)
			
			hasScript := content != ""
			if hasScript != tt.expectedHasScript {
				t.Errorf("extractScriptSection() hasScript = %v, want %v", 
					hasScript, tt.expectedHasScript)
			}
			
			if startLine != tt.expectedStartLine {
				t.Errorf("extractScriptSection() startLine = %d, want %d", 
					startLine, tt.expectedStartLine)
			}
		})
	}
}

func TestIsHTMLTag(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		expected bool
	}{
		{"div", "div", true},
		{"span", "span", true},
		{"button", "button", true},
		{"form", "form", true},
		{"input", "input", true},
		{"uppercase div", "DIV", false},
		{"mixed case", "Div", false},
		{"component", "MyComponent", false},
		{"kebab component", "my-component", false},
		{"q-form", "q-form", false},
		{"QForm", "QForm", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHTMLTag(tt.tagName)
			if result != tt.expected {
				t.Errorf("isHTMLTag(%q) = %v, want %v", tt.tagName, result, tt.expected)
			}
		})
	}
}
