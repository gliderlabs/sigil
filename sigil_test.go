package sigil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLookPath(t *testing.T) {
	for _, tmpDir := range []string{"", "."} {
		for _, relative := range []bool{true, false} {
			for _, exists := range []bool{true, false} {
				// setup
				f, err := os.CreateTemp(tmpDir, "sigiltest")
				if err != nil {
					t.Error("failed to crate temporary file")
					continue
				}
				relpath := f.Name()
				fullpath, err := filepath.Abs(f.Name())
				if err != nil {
					t.Error("failed to get the absolute path")
				}
				if !exists {
					os.Remove(f.Name())
				}
				path := fullpath
				if relative {
					path = relpath
				}

				// test
				p, err := LookPath(path)
				if exists {
					if p != fullpath {
						t.Errorf("expected %s but %s; tmpDir=%s, relative=%v, exists=%v", fullpath, p, tmpDir, relative, exists)
					}
				} else {
					if err == nil {
						t.Errorf("expected error. tmpDir=%s, relative=%v, exists=%v", tmpDir, relative, exists)
					}
				}

				// teardown
				f.Close()
				os.Remove(f.Name())
			}
		}
	}
}

func writeTempFile(t *testing.T, ext, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "sigil-test-*"+ext)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestParseVarsFileJSON(t *testing.T) {
	path := writeTempFile(t, ".json", `{"name": "Jeff", "greeting": "Hello"}`)
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", vars["name"])
	}
	if vars["greeting"] != "Hello" {
		t.Errorf("expected greeting=Hello, got %v", vars["greeting"])
	}
}

func TestParseVarsFileJSONNested(t *testing.T) {
	path := writeTempFile(t, ".json", `{"db": {"host": "localhost", "port": 5432}}`)
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	db, ok := vars["db"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected db to be map, got %T", vars["db"])
	}
	if db["host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", db["host"])
	}
}

func TestParseVarsFileYAML(t *testing.T) {
	path := writeTempFile(t, ".yaml", "name: Jeff\ngreeting: Hello\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", vars["name"])
	}
	if vars["greeting"] != "Hello" {
		t.Errorf("expected greeting=Hello, got %v", vars["greeting"])
	}
}

func TestParseVarsFileYML(t *testing.T) {
	path := writeTempFile(t, ".yml", "name: Jeff\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", vars["name"])
	}
}

func TestParseVarsFileYAMLNested(t *testing.T) {
	path := writeTempFile(t, ".yaml", "db:\n  host: localhost\n  port: 5432\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	db, ok := vars["db"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected db to be map, got %T", vars["db"])
	}
	if db["host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", db["host"])
	}
}

func TestParseVarsFileEnv(t *testing.T) {
	content := "# comment\nname=Jeff\ngreeting=\"Hello World\"\n\nexport FOO='bar'\n"
	path := writeTempFile(t, ".env", content)
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", vars["name"])
	}
	if vars["greeting"] != "Hello World" {
		t.Errorf("expected greeting='Hello World', got %v", vars["greeting"])
	}
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", vars["FOO"])
	}
}

func TestParseVarsFileEnvValueWithEquals(t *testing.T) {
	path := writeTempFile(t, ".env", "connection=host=localhost dbname=test\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["connection"] != "host=localhost dbname=test" {
		t.Errorf("expected value with equals preserved, got %v", vars["connection"])
	}
}

func TestParseVarsFileEnvEmptyValue(t *testing.T) {
	path := writeTempFile(t, ".env", "empty=\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["empty"] != "" {
		t.Errorf("expected empty value, got %v", vars["empty"])
	}
}

func TestParseVarsFileUnknownExtension(t *testing.T) {
	path := writeTempFile(t, ".txt", "name=Jeff\ncolor=red\n")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", vars["name"])
	}
	if vars["color"] != "red" {
		t.Errorf("expected color=red, got %v", vars["color"])
	}
}

func TestParseVarsFileNotFound(t *testing.T) {
	_, err := ParseVarsFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseVarsFileInvalidJSON(t *testing.T) {
	path := writeTempFile(t, ".json", "not json")
	defer os.Remove(path)

	_, err := ParseVarsFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseVarsFileInvalidYAML(t *testing.T) {
	path := writeTempFile(t, ".yaml", "- one\n- two\n")
	defer os.Remove(path)

	_, err := ParseVarsFile(path)
	if err == nil {
		t.Error("expected error for YAML list at top level")
	}
}

func TestParseVarsFileEmptyYAML(t *testing.T) {
	path := writeTempFile(t, ".yaml", "")
	defer os.Remove(path)

	vars, err := ParseVarsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(vars) != 0 {
		t.Errorf("expected empty map, got %v", vars)
	}
}

func TestParseEnvContentExportWithoutSpace(t *testing.T) {
	vars, err := parseEnvContent([]byte("exportFOO=bar\n"))
	if err != nil {
		t.Fatal(err)
	}
	if vars["exportFOO"] != "bar" {
		t.Errorf("expected exportFOO=bar, got %v", vars)
	}
}

func TestParseEnvContentRejectsNoEquals(t *testing.T) {
	_, err := parseEnvContent([]byte("noequals\nname=Jeff\n"))
	if err == nil {
		t.Error("expected error for line without =")
	}
}

func TestConvertYAMLMap(t *testing.T) {
	input := map[interface{}]interface{}{
		"name": "Jeff",
		"nested": map[interface{}]interface{}{
			"key": "value",
		},
	}
	result := convertYAMLMap(input)
	if result["name"] != "Jeff" {
		t.Errorf("expected name=Jeff, got %v", result["name"])
	}
	nested, ok := result["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested to be map[string]interface{}, got %T", result["nested"])
	}
	if nested["key"] != "value" {
		t.Errorf("expected nested.key=value, got %v", nested["key"])
	}
}

func TestConvertYAMLValueSlice(t *testing.T) {
	input := []interface{}{
		map[interface{}]interface{}{"a": "b"},
		"plain",
	}
	result := convertYAMLValue(input)
	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}
	m, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map in slice, got %T", slice[0])
	}
	if !reflect.DeepEqual(m, map[string]interface{}{"a": "b"}) {
		t.Errorf("unexpected map: %v", m)
	}
}
