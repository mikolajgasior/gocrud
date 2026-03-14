package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
)

// HTMLRenderer implements the Renderer interface for HTML output
type HTMLRenderer struct {
	template  *template.Template
	templates map[string]*template.Template
}

// RegisterTemplate allows registering custom templates for different types
func (h *HTMLRenderer) RegisterTemplate(name string, tmplStr string) error {
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return err
	}
	h.templates[name] = tmpl
	return nil
}

// Render converts any object to HTML string using templates.
// It does NOT fall back to generic rendering. If a type is unknown and not registered,
// it returns an error message.
func (h *HTMLRenderer) Render(obj interface{}) string {
	if obj == nil {
		return "<!-- null object -->"
	}

	// Check if object implements Renderable interface
	if renderable, ok := obj.(Renderable); ok {
		return renderable.RenderString()
	}

	// Get the type of the object
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	// Handle pointer types
	if objType.Kind() == reflect.Ptr {
		if objValue.IsNil() {
			return "<!-- nil pointer -->"
		}
		objType = objType.Elem()
		objValue = objValue.Elem()
	}

	// Handle Input struct specifically (built-in support)
	if objType.Name() == "Input" {
		return h.renderWithTemplate(objValue.Interface(), "Input")
	}

	// For other structs, look for a registered template
	templateName := objType.Name()
	if tmpl, exists := h.templates[templateName]; exists {
		return h.executeTemplate(tmpl, objValue.Interface())
	}

	// NO GENERIC FALLBACK
	// If we reach here, the type is unknown and no template was registered.
	return fmt.Sprintf("<!-- Error: No rendering strategy found for type '%s'. Register a template or implement Renderable. -->", templateName)
}

// renderWithTemplate executes a named template with the given object
func (h *HTMLRenderer) renderWithTemplate(obj interface{}, templateName string) string {
	tmpl, exists := h.templates[templateName]
	if !exists {
		return fmt.Sprintf("<!-- template '%s' not found -->", templateName)
	}
	return h.executeTemplate(tmpl, obj)
}

// executeTemplate executes a template and returns the result
func (h *HTMLRenderer) executeTemplate(tmpl *template.Template, obj interface{}) string {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, obj)
	if err != nil {
		return fmt.Sprintf("<!-- render error: %v -->", err)
	}
	return buf.String()
}
