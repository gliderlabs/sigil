package fs

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/spf13/afero"
)

type fsTest struct {
	args        []interface{}
	expected    interface{}
	expectedErr interface{}
}

func TestModuleFunc(t *testing.T) {
	t.Parallel()
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "file.txt", []byte("hello"), 0644)
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"fs": ModuleFunc(fs),
	}).Parse(`{{ fs.Read "file.txt" }}`)
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
