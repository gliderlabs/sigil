package str

import (
	"bytes"
	"html/template"
	"testing"
)

type strTest struct {
	args        []interface{}
	expected    interface{}
	expectedErr error
}

func TestModuleFunc(t *testing.T) {
	t.Parallel()
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"str": ModuleFunc,
	}).Parse(`{{ "Hello" | str.Upper }}`)
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Error(err)
	}
	if got := buf.String(); got != "HELLO" {
		t.Errorf("unexpected executed template content: %s", got)
	}
}
