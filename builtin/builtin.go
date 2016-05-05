package builtin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/dustin/go-jsonpointer"
	"github.com/flynn/go-shlex"
	yyaml "github.com/ghodss/yaml"
	"github.com/gliderlabs/sigil"
	"gopkg.in/yaml.v2"
)

func init() {
	sigil.Register(template.FuncMap{
		// templating
		"include": Include,
		"default": Default,
		"var":     Var,
		// strings
		"capitalize": Capitalize,
		"lower":      Lower,
		"upper":      Upper,
		"replace":    Replace,
		"trim":       Trim,
		"indent":     Indent,
		"match":      Match,
		"render":     Render,
		"stdin":      Stdin,
		"substr":     Substring,
		// filesystem
		"file":   File,
		"exists": Exists,
		"dir":    Dir,
		"dirs":   Dirs,
		"files":  Files,
		"text":   Text,
		// external
		"sh":      Shell,
		"httpget": HttpGet,
		// structured data
		"pointer":    Pointer,
		"json":       Json,
		"tojson":     ToJson,
		"yaml":       Yaml,
		"toyaml":     ToYaml,
		"uniq":       Uniq,
		"drop":       Drop,
		"append":     Append,
		"seq":        Seq,
		"join":       Join,
		"joinkv":     JoinKv,
		"split":      Split,
		"splitkv":    SplitKv,
		"yamltojson": YamlToJson,
	})
}

func Shell(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("sh must be given a string")
	}
	args, err := shlex.Split(in_)
	if err != nil {
		return "", err
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	out, err := exec.Command(path, args[1:]...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func HttpGet(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("httpget must be given a string")
	}
	resp, err := http.Get(in_)
	if err != nil {
		return "", err
	}
	return sigil.NamedReader{resp.Body, "<"+in_+">"}, nil
}


func JoinKv(sep string, in interface{}) ([]interface{}, error) {
	m, ok := in.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("joinkv must be given a string map of strings")
	}
	var elements []interface{}
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("joinkv must be given a string map of strings")
		}
		elements = append(elements, strings.Join([]string{k, s}, sep))
	}
	return elements, nil
}

func SplitKv(sep string, in []interface{}) (interface{}, error) {
	out := make(map[string]interface{})
	for i := range in {
		v, ok := in[i].(string)
		if !ok {
			return nil, fmt.Errorf("joinkv must be given a string map of strings")
		}
		parts := strings.SplitN(v, sep, 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		} else {
			out[v] = true
		}
	}
	return out, nil
}

func Seq(i interface{}) ([]interface{}, error) {
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
	var el []interface{}
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

func Join(delim string, in []interface{}) interface{} {
	var elements []string
	for _, el := range in {
		str, ok := el.(string)
		if ok {
			elements = append(elements, str)
		}
	}
	return strings.Join(elements, delim)
}

func Split(delim string, in interface{}) ([]interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("split must be given a string")
	}
	var elements []interface{}
	for _, v := range strings.Split(in_, delim) {
		elements = append(elements, v)
	}
	return elements, nil
}

func Substring(slice string, in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("substr must be given a string")
	}
	s := strings.Split(slice, ":")
	start, err := strconv.Atoi(s[0])
	if err != nil {
		start = 0
	}
	end, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, fmt.Errorf("substr needs slice expression as 'start:end' ")
	}
	return in_[start:end], nil
}

func Capitalize(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("capitalize must be given a string")
	}
	return strings.Title(in_), nil
}

func Lower(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("lower must be given a string")
	}
	return strings.ToLower(in_), nil
}

func Upper(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("upper must be given a string")
	}
	return strings.ToUpper(in_), nil
}

func Replace(old, new string, in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("replace must be given a string")
	}
	return strings.Replace(in_, old, new, -1), nil
}

func Trim(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("trim must be given a string")
	}
	return strings.Trim(in_, " \n"), nil
}

func read(file interface{}) ([]byte, error) {
	reader, ok := file.(sigil.NamedReader)
	if ok {
		return ioutil.ReadAll(reader)
	}
	path, _, ok := sigil.String(file)
	if !ok {
		return []byte{}, fmt.Errorf("file must be stream or path string")
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

func File(filename interface{}) (interface{}, error) {
	str, err := read(filename)
	return string(str), err
}

func Text(file interface{}) (interface{}, error) {
	f, err := read(file)
	if err != nil {
		return nil, err
	}
	return string(f), nil
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

func ToJson(obj interface{}) (interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return string(data), nil
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

func ToYaml(obj interface{}) (interface{}, error) {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func YamlToJson(in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("ymltojson must be given a string")
	}

	j2, err := yyaml.YAMLToJSON([]byte(in_))
	if err != nil {
		return nil, err
	}
	return string(j2), nil
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

func Render(args ...interface{}) (interface{}, error) {
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
	vars := make(map[string]interface{})
	for _, arg := range args {
		mv, ok := arg.(map[string]interface{})
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

func Include(filename string, args ...interface{}) (interface{}, error) {
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

func Indent(indent string, in interface{}) (interface{}, error) {
	in_, _, ok := sigil.String(in)
	if !ok {
		return "", fmt.Errorf("indent must be given a string")
	}
	var indented []string
	lines := strings.Split(in_, "\n")
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
	return strings.Join(indented, "\n"), nil
}

func Var(name string) interface{} {
	return os.Getenv(name)
}

func Match(pattern string, in interface{}) (bool, error) {
	str, _, ok := sigil.String(in)
	if !ok {
		return false, fmt.Errorf("match must be given a string")
	}
	return path.Match(pattern, str)
}

func Exists(in interface{}) (bool, error) {
	filename, _, ok := sigil.String(in)
	if !ok {
		return false, fmt.Errorf("exists must be given a string")
	}
	_, err := sigil.LookPath(filename)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func Dir(in interface{}) ([]interface{}, error) {
	path, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("dir must be given a string")
	}
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

func Dirs(in interface{}) ([]interface{}, error) {
	path, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("dirs must be given a string")
	}
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

func Files(in interface{}) ([]interface{}, error) {
	path, _, ok := sigil.String(in)
	if !ok {
		return nil, fmt.Errorf("files must be given a string")
	}
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

func Stdin() (interface{}, error) {
	return sigil.NamedReader{os.Stdin, "<stdin>"}, nil
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
