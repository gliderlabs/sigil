package builtin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func file(filename string) []byte {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = filepath.Join(sigil.TemplateDir, filename)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return []byte{}
		}
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return data
}

func File(filename string) string {
	return string(file(filename))
}

func Json(filename string) map[string]interface{} {
	var obj map[string]interface{}
	err := json.Unmarshal(file(filename), &obj)
	if err != nil {
		return nil
	}
	return obj
}

func Yaml(filename string) map[string]interface{} {
	var obj map[string]interface{}
	err := yaml.Unmarshal(file(filename), &obj)
	if err != nil {
		return nil
	}
	return obj
}

func Pointer(path string, in map[string]interface{}) interface{} {
	return jsonpointer.Get(in, path)
}

func Include(filename string, args ...string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	vars := make(map[string]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	str, err := sigil.Execute(string(data), vars)
	if err != nil {
		return "", err
	}
	return str, nil
}
