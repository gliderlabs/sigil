package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gliderlabs/sigil"
	_ "github.com/gliderlabs/sigil/builtin"
	flag "github.com/spf13/pflag"
)

var Version string

var (
	filename  = flag.StringP("filename", "f", "", "use template file instead of STDIN")
	inline    = flag.StringP("inline", "i", "", "use inline template string instead of STDIN")
	inPlace   = flag.Bool("in-place", false, "write output back to the file specified by -f")
	varsFiles = flag.StringArrayP("vars-file", "V", []string{}, "load variables from a file (JSON, YAML, or env format)")
	posix     = flag.BoolP("posix", "p", false, "preprocess with POSIX variable expansion")
	version   = flag.BoolP("version", "v", false, "prints version")
)

func template() ([]byte, string, error) {
	if *filename != "" {
		data, err := os.ReadFile(*filename)
		if err != nil {
			return []byte{}, "", err
		}
		sigil.PushPath(filepath.Dir(*filename))
		return data, filepath.Base(*filename), nil
	}
	if *inline != "" {
		return []byte(*inline), "<inline>", nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return []byte{}, "", err
	}
	return data, "<stdin>", nil
}

func writeInPlace(filename string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(filename)
	tmp, err := os.CreateTemp(dir, ".sigil-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmp.Name()

	success := false
	defer func() {
		if !success {
			os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}
	if err := os.Rename(tmpName, filename); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}
	success = true
	return nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *inPlace && *filename == "" {
		fmt.Fprintln(os.Stderr, "--in-place requires -f/--filename")
		os.Exit(1)
	}
	if *posix {
		sigil.PosixPreprocess = true
	}
	if os.Getenv("SIGIL_PATH") != "" {
		sigil.TemplatePath = strings.Split(os.Getenv("SIGIL_PATH"), ":")
	}

	var originalMode os.FileMode
	if *inPlace {
		info, err := os.Stat(*filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		originalMode = info.Mode()
	}

	vars := make(map[string]interface{})
	for _, vf := range *varsFiles {
		fileVars, err := sigil.ParseVarsFile(vf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for k, v := range fileVars {
			vars[k] = v
		}
	}
	for _, arg := range flag.Args() {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	tmpl, name, err := template()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	render, err := sigil.Execute(tmpl, vars, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *inPlace {
		if err := writeInPlace(*filename, render.Bytes(), originalMode); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		os.Stdout.Write(render.Bytes())
	}
}
