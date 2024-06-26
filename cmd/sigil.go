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
	filename = flag.StringP("filename", "f", "", "use template file instead of STDIN")
	inline   = flag.StringP("inline", "i", "", "use inline template string instead of STDIN")
	posix    = flag.BoolP("posix", "p", false, "preprocess with POSIX variable expansion")
	version  = flag.BoolP("version", "v", false, "prints version")
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

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *posix {
		sigil.PosixPreprocess = true
	}
	if os.Getenv("SIGIL_PATH") != "" {
		sigil.TemplatePath = strings.Split(os.Getenv("SIGIL_PATH"), ":")
	}
	vars := make(map[string]interface{})
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
	os.Stdout.Write(render.Bytes())
}
