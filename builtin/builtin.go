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
		"exists":     Exists,
		"dir":        Dir,
		"dirs":       Dirs,
		"files":      Files,
		"uniq":       Uniq,
		"drop":       Drop,
		"append":     Append,
		"stdin":      Stdin,
		"entries":    Entries,
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
	if in == nil {
		return value
	}
	if reflect.Zero(reflect.TypeOf(in)).Interface() == in {
		return value
	}
	return in
}

func Join(delim string, in []interface{}) string {
	var elements []string
	for _, el := range in {
		str, ok := el.(string)
		if ok {
			elements = append(elements, str)
		}
	}
	return strings.Join(elements, delim)
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

func read(file interface{}) ([]byte, error) {
	stdin, ok := file.(stdinStr)
	if ok {
		return []byte(stdin), nil
	}
	path, ok := file.(string)
	if !ok {
		return []byte{}, fmt.Errorf("file must be path string or stdin")
	}
	filepath, err := sigil.LookPath(path)
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
	str, err := read(filename)
	return string(str), err
}

func Json(file interface{}) (interface{}, error) {
	var obj interface{}
	f, err := read(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func Yaml(file interface{}) (interface{}, error) {
	var obj interface{}
	f, err := read(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func Pointer(path string, in interface{}) (interface{}, error) {
	m := make(map[string]interface{})
	switch val := in.(type) {
	case map[string]interface{}:
		for k, v := range val {
			m[k] = v
		}
	case map[interface{}]interface{}:
		for k, v := range val {
			m[k.(string)] = v
		}
	default:
		return nil, fmt.Errorf("pointer needs a map type")
	}
	return jsonpointer.Get(m, path), nil
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
			if line != "" {
				indented = append(indented, indent+line)
			} else {
				indented = append(indented, line)
			}
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

func Exists(filename string) bool {
	_, err := sigil.LookPath(filename)
	if err != nil {
		return false
	}
	return true
}

func Dir(path string) ([]interface{}, error) {
	var files []interface{}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		files = append(files, fi.Name())
	}
	return files, nil
}

func Dirs(path string) ([]interface{}, error) {
	var dirs []interface{}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		}
	}
	return dirs, nil
}

func Files(path string) ([]interface{}, error) {
	var files []interface{}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if !fi.IsDir() {
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

func Uniq(in ...[]interface{}) []interface{} {
	m := make(map[interface{}]bool)
	for i := range in {
		for _, v := range in[i] {
			m[v] = true
		}
	}
	var uniq []interface{}
	for k, _ := range m {
		uniq = append(uniq, k)
	}
	return uniq
}

type stdinStr string

func Stdin() (stdinStr, error) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return stdinStr(data), nil
}

func Append(item interface{}, items []interface{}) []interface{} {
	return append(items, item)
}

func Drop(item interface{}, items []interface{}) ([]interface{}, error) {
	var out []interface{}
	pattern, isstr := item.(string)
	if isstr {
		for i := range items {
			str, ok := items[i].(string)
			if !ok {
				return nil, fmt.Errorf("all elements must be a string to drop a string")
			}
			match, err := path.Match(pattern, str)
			if err != nil {
				return nil, fmt.Errorf("bad pattern: %s", pattern)
			}
			if !match {
				out = append(out, items[i])
			}
		}
		return out, nil
	}
	for i := range items {
		if item != items[i] {
			out = append(out, items[i])
		}
	}
	return out, nil
}

func Entries(pattern string, m map[interface{}]interface{}) []interface{} {
	ret := []interface{}{}
	for k, v := range m {
		ret = append(ret, fmt.Sprintf(pattern, k, v))
	}
	return ret
}
