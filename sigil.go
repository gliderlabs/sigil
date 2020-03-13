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

var leftDelim = "{{"
var rightDelim = "}}"
var fnMap = template.FuncMap{}

func init() {
	delims := os.Getenv("SIGIL_DELIMS")
	if delims != "" {
		d := strings.Split(delims, ",")
		leftDelim = d[0]
		rightDelim = d[1]
	}
}

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
		if fileExists(file) {
			return file, nil
		}
	} else {
		cwd, _ := os.Getwd()
		search := append([]string{cwd}, TemplatePath...)
		for _, path := range search {
			filepath := filepath.Join(path, file)
			if fileExists(filepath) {
				return filepath, nil
			}
		}
	}
	return "", fmt.Errorf("Not found in path: %s %v", file, TemplatePath)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func restoreEnv(env []string) {
	os.Clearenv()
	for _, kvp := range env {
		kv := strings.SplitN(kvp, "=", 2)
		os.Setenv(kv[0], kv[1])
	}
}

func Execute(input []byte, vars map[string]interface{}, name string) (bytes.Buffer, error) {
	var tmplVars string
	var err error
	defer restoreEnv(os.Environ())
	for k, iv := range vars {
		if v, ok := iv.(string); ok {
			err := os.Setenv(k, v)
			if err != nil {
				return bytes.Buffer{}, err
			}
		}
		tmplVars = tmplVars + fmt.Sprintf("%s $%s := .%s %s", leftDelim, k, k, rightDelim)
	}
	inputStr := string(input)
	if PosixPreprocess {
		inputStr, err = posix.ExpandEnv(inputStr)
		if err != nil {
			return bytes.Buffer{}, err
		}
	}

	inputStr = strings.Replace(
		inputStr,
		fmt.Sprintf("\\%s\n%s", rightDelim, leftDelim),
		fmt.Sprintf("%s%s", rightDelim, leftDelim),
		-1,
	)
	tmpl, err := template.New(name).Funcs(fnMap).Delims(leftDelim, rightDelim).Parse(tmplVars + inputStr)
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
