package sigil

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mgood/go-posix"
)

var (
	TemplatePath    []string
	PosixPreprocess bool
)

var fnMap = template.FuncMap{}

func Register(fm template.FuncMap) {
	for k, v := range fm {
		fnMap[k] = v
	}
}

func PushPath(path string) {
	TemplatePath = append([]string{path}, TemplatePath...)
}

func PopPath() {
	_, TemplatePath = TemplatePath[0], TemplatePath[1:]
}

func LookPath(file string) (string, error) {
	cwd, _ := os.Getwd()
	search := append([]string{cwd}, TemplatePath...)
	for _, path := range search {
		filepath := filepath.Join(path, file)
		if _, err := os.Stat(filepath); err == nil {
			return filepath, nil
		}
	}
	return "", fmt.Errorf("Not found in path: %s %v", file, TemplatePath)
}

func Execute(input string, vars map[string]string, name string) (string, error) {
	var tmplVars string
	var err error
	for k, v := range vars {
		err := os.Setenv(k, v)
		if err != nil {
			return "", err
		}
		escaped := strings.Replace(v, "\"", "\\\"", -1)
		tmplVars = tmplVars + fmt.Sprintf("{{ $%s := \"%s\" }}", k, escaped)
	}
	if PosixPreprocess {
		input, err = posix.ExpandEnv(input)
		if err != nil {
			return "", err
		}
	}
	input = strings.Replace(input, "}}\n{{", "}}{{", -1)
	tmpl, err := template.New(name).Funcs(fnMap).Parse(tmplVars + input)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
