package var_

import (
	"bytes"
	"html/template"
	"testing"
)

type varTest struct {
	args        []interface{}
	expected    interface{}
	expectedErr error
}

func TestModuleFunc(t *testing.T) {
	t.Parallel()
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"var": ModuleFunc,
	}).Parse(`{{ "" | var.Default "hello" }}`)
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Error(err)
	}
	if got := buf.String(); got != "hello" {
		t.Errorf("unexpected executed template content: %s", got)
	}
}
