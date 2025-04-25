package generator

import (
	"bytes"
	"testing"
	"text/template"
)

// TestBasicTemplateRendering tests that the template engine properly renders templates
func TestBasicTemplateRendering(t *testing.T) {
	templateStr := "Hello, {{ .Name }}!"
	data := map[string]string{
		"Name": "World",
	}

	tmpl, err := template.New("test").Parse(templateStr)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := "Hello, World!"
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

// TestTemplateVariableInterpolation tests that variables are correctly interpolated
func TestTemplateVariableInterpolation(t *testing.T) {
	templateStr := `package {{ .PackageName }}

import (
	"{{ .ImportPath }}"
)

// {{ .StructName }} represents a {{ .Description }}
type {{ .StructName }} struct {
	{{ range .Fields }}{{ .Name }} {{ .Type }}
	{{ end }}
}
`

	data := map[string]interface{}{
		"PackageName": "models",
		"ImportPath":  "github.com/example/pkg/types",
		"StructName":  "User",
		"Description": "user entity",
		"Fields": []map[string]string{
			{"Name": "ID", "Type": "string"},
			{"Name": "Name", "Type": "string"},
			{"Name": "Email", "Type": "string"},
		},
	}

	tmpl, err := template.New("model").Parse(templateStr)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Check that essential elements are in the output
	output := buf.String()
	expectedParts := []string{
		"package models",
		"import (",
		"\"github.com/example/pkg/types\"",
		"type User struct",
		"ID string",
		"Name string",
		"Email string",
	}

	for _, part := range expectedParts {
		if !bytes.Contains(buf.Bytes(), []byte(part)) {
			t.Errorf("Expected output to contain %q but it didn't.\nOutput: %s", part, output)
		}
	}
}
