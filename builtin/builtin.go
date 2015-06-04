package builtin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/dustin/go-jsonpointer"
	"github.com/gliderlabs/sigil"
	"gopkg.in/yaml.v2"
)

func init() {
	sigil.Register(template.FuncMap{
		"seq":        Seq,
		"default":    Default,
		"join":       Join,
		"split":      Split,
		"capitalize": Capitalize,
		"lower":      Lower,
		"upper":      Upper,
		"replace":    Replace,
		"trim":       Trim,
		"file":       File,
		"json":       Json,
		"yaml":       Yaml,
		"pointer":    Pointer,
		"include":    Include,
		"indent":     Indent,
		"var":        Var,
		"match":      Match,
		"render":     Render,
	})
}

func Seq(i interface{}) ([]string, error) {
	var num int
	var err error
	var valid bool
	switch v := i.(type) {
	case int, int32, int64:
		num, valid = v.(int)
	case string:
		num, err = strconv.Atoi(v)
		if err == nil {
			valid = true
		}
	}
	if !valid {
		return nil, fmt.Errorf("seq must be given an integer or numeric string")
	}
	var el []string
	for i, _ := range make([]bool, num) {
		el = append(el, strconv.Itoa(i))
	}
	return el, nil
}

func Default(value, in interface{}) interface{} {
	if reflect.Zero(reflect.TypeOf(in)).Interface() == in {
		return value
	}
	return in
}

func Join(delim string, in []string) string {
	return strings.Join(in, delim)
}

func Split(delim string, in string) []string {
	return strings.Split(in, delim)
}

func Capitalize(in string) string {
	return strings.Title(in)
}

func Lower(in string) string {
	return strings.ToLower(in)
}

func Upper(in string) string {
	return strings.ToUpper(in)
}

func Replace(old, new, in string) string {
	return strings.Replace(in, old, new, -1)
}

func Trim(in string) string {
	return strings.Trim(in, " \n")
}

func file(file string) ([]byte, error) {
	filepath, err := sigil.LookPath(file)
	if err != nil {
		return []byte{}, err
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func File(filename string) (string, error) {
	str, err := file(filename)
	return string(str), err
}

func Json(filename string) (interface{}, error) {
	var obj interface{}
	f, err := file(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func Yaml(filename string) (interface{}, error) {
	var obj interface{}
	f, err := file(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func Pointer(path string, in map[interface{}]interface{}) interface{} {
	m := make(map[string]interface{})
	for k, v := range in {
		m[k.(string)] = v
	}
	return jsonpointer.Get(m, path)
}

func Render(args ...interface{}) (string, error) {
	if len(args) == 0 {
		fmt.Errorf("render cannot be used without arguments")
	}
	input := args[len(args)-1].(string)
	var vars []interface{}
	if len(args) > 1 {
		vars = args[0 : len(args)-1]
	}
	render, err := render([]byte(input), vars, "<render>")
	return render.String(), err
}

func render(data []byte, args []interface{}, name string) (bytes.Buffer, error) {
	vars := make(map[string]string)
	for _, arg := range args {
		mv, ok := arg.(map[string]string)
		if ok {
			for k, v := range mv {
				vars[k] = v
			}
			continue
		}
		sv, ok := arg.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(sv, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return sigil.Execute(data, vars, name)
}

func Include(filename string, args ...interface{}) (string, error) {
	path, err := sigil.LookPath(filename)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	sigil.PushPath(filepath.Dir(path))
	defer sigil.PopPath()
	render, err := render(data, args, filepath.Base(path))
	return render.String(), err
}

func Indent(indent, in string) string {
	var indented []string
	lines := strings.Split(in, "\n")
	indented = append(indented, lines[0])
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			indented = append(indented, indent+line)
		}
	}
	return strings.Join(indented, "\n")
}

func Var(name string) string {
	return os.Getenv(name)
}

func Match(pattern string, str string) (bool, error) {
	return path.Match(pattern, str)
}
