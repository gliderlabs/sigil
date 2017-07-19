package b64

import (
	"bytes"
	"html/template"
	"testing"
)

type b64Test struct {
	args        []interface{}
	expected    interface{}
	expectedErr error
}

func TestModuleFunc(t *testing.T) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"b64": ModuleFunc,
	}).Parse(`{{ "Hello" | b64.Encode }}`)
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Error(err)
	}
	if got := buf.String(); got != "SGVsbG8=" {
		t.Errorf("unexpected executed template content: %s", got)
	}
}
