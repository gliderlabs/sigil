package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gliderlabs/sigil"
	_ "github.com/gliderlabs/sigil/builtin"
)

var Version string

var (
	filename = flag.String("f", "", "use template file instead of STDIN")
	version  = flag.Bool("v", false, "prints version")
)

func template() (string, error) {
	if *filename != "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			return "", err
		}
		sigil.TemplateDir = filepath.Dir(*filename)
		return string(data), nil
	}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	vars := make(map[string]string)
	for _, arg := range os.Args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	tmpl, err := template()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	render, err := sigil.Execute(tmpl, vars)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(render)
}
