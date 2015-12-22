package sigil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

type NamedReader struct {
	io.Reader
	Name string
}

func String(in interface{}) (string, string, bool) {
	switch obj := in.(type) {
	case string:
		return obj, "", true
	case NamedReader:
		data, err := ioutil.ReadAll(obj)
		if err != nil {
			// TODO: better overall error/panic handling
			panic(err)
		}
		return string(data), obj.Name, true
	case fmt.Stringer:
		return obj.String(), "", true
	default:
		return "", "", false
	}
}

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
	if strings.HasPrefix(file, "/") {
		return file, nil
	}
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

func restoreEnv(env []string) {
	os.Clearenv()
	for _, kvp := range env {
		kv := strings.SplitN(kvp, "=", 2)
		os.Setenv(kv[0], kv[1])
	}
}

func Execute(input []byte, vars map[string]string, name string) (bytes.Buffer, error) {
	var tmplVars string
	var err error
	defer restoreEnv(os.Environ())
	for k, v := range vars {
		err := os.Setenv(k, v)
		if err != nil {
			return bytes.Buffer{}, err
		}
		escaped := strings.Replace(v, "\"", "\\\"", -1)
		tmplVars = tmplVars + fmt.Sprintf("{{ $%s := \"%s\" }}", k, escaped)
	}
	inputStr := string(input)
	if PosixPreprocess {
		inputStr, err = posix.ExpandEnv(inputStr)
		if err != nil {
			return bytes.Buffer{}, err
		}
	}
	inputStr = strings.Replace(inputStr, "\\}}\n{{", "}}{{", -1)
	tmpl, err := template.New(name).Funcs(fnMap).Parse(tmplVars + inputStr)
	if err != nil {
		return bytes.Buffer{}, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}
