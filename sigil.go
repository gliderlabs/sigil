package sigil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mgood/go-posix"
	"gopkg.in/yaml.v2"
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
		data, err := io.ReadAll(obj)
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

func convertYAMLValue(v interface{}) interface{} {
	switch t := v.(type) {
	case map[interface{}]interface{}:
		return convertYAMLMap(t)
	case []interface{}:
		for i, elem := range t {
			t[i] = convertYAMLValue(elem)
		}
		return t
	default:
		return v
	}
}

func convertYAMLMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[fmt.Sprintf("%v", k)] = convertYAMLValue(v)
	}
	return result
}

func parseEnvContent(data []byte) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}
		vars[key] = value
	}
	return vars, nil
}

// ParseVarsFile reads a file and parses it as template variables.
// The format is auto-detected by file extension:
// .json files are parsed as JSON objects, .yaml/.yml as YAML mappings,
// and all other extensions as key=value lines.
func ParseVarsFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vars file %s: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		var obj map[string]interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("failed to parse JSON vars file %s: %w", path, err)
		}
		return obj, nil
	case ".yaml", ".yml":
		var obj interface{}
		if err := yaml.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("failed to parse YAML vars file %s: %w", path, err)
		}
		if obj == nil {
			return make(map[string]interface{}), nil
		}
		m, ok := obj.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("YAML vars file %s must contain a mapping at top level", path)
		}
		return convertYAMLMap(m), nil
	default:
		return parseEnvContent(data)
	}
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
