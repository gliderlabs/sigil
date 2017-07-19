package main

import (
	"flag"
	f "fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gliderlabs/sigil/pkg/sigil"
	"github.com/spf13/afero"

	"github.com/gliderlabs/sigil/lib/b64"
	"github.com/gliderlabs/sigil/lib/fmt"
	"github.com/gliderlabs/sigil/lib/fs"
	"github.com/gliderlabs/sigil/lib/str"
	"github.com/gliderlabs/sigil/lib/var"
)

var Version string

var (
	filename = flag.String("f", "", "use template file instead of STDIN")
	inline   = flag.String("i", "", "use inline template string instead of STDIN")
	posix    = flag.Bool("P", false, "use POSIX variable expansion mode")
	version  = flag.Bool("v", false, "prints version")
)

func NewProcessor() (*sigil.Processor, error) {
	var Template string
	if *filename != "" && *inline == "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			return nil, err
		}
		Template = string(data)
	}
	if *inline != "" && *filename == "" {
		Template = *inline
	}
	if Template == "" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		Template = string(data)
	}
	return &sigil.Processor{
		PosixMode: *posix,
		Template:  Template,
	}, nil
}

func ParseVars() map[string]interface{} {
	vars := make(map[string]interface{})
	for _, arg := range flag.Args() {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}

func main() {
	flag.Parse()
	if *version {
		if Version == "" {
			Version = "master"
		}
		f.Println(Version)
		os.Exit(0)
	}
	p, err := NewProcessor()
	if err != nil {
		f.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	p.RegisterFuncs(template.FuncMap{
		"b64": b64.ModuleFunc,
		"fmt": fmt.ModuleFunc,
		"fs":  fs.ModuleFunc(afero.NewOsFs()),
		"str": str.ModuleFunc,
		"var": var_.ModuleFunc,
	})
	if err := p.Execute(os.Stdout, ParseVars()); err != nil {
		f.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
