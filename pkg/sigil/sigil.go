package sigil

import (
	"html/template"
	"io"
	"os"
	"strings"

	posix "github.com/mgood/go-posix"
)

type Processor struct {
	PosixMode bool
	Paths     []string
	Template  string

	funcMap    template.FuncMap
	delimLeft  string
	delimRight string
}

func (p *Processor) RegisterFuncs(f template.FuncMap) {
	if p.funcMap == nil {
		p.funcMap = make(template.FuncMap)
	}
	for k, v := range f {
		p.funcMap[k] = v
	}
}

func (p *Processor) PushPath(path string) {
	p.Paths = append([]string{path}, p.Paths...)
}

func (p *Processor) PopPath() {
	_, p.Paths = p.Paths[0], p.Paths[1:]
}

func (p *Processor) Execute(w io.Writer, vars map[string]interface{}) error {
	if p.PosixMode {
		return p.executePosix(w, vars)
	} else {
		return p.executeTmpl(w, vars)
	}
}

func (p *Processor) executePosix(w io.Writer, vars map[string]interface{}) error {
	defer func(env []string) {
		os.Clearenv()
		for _, kvp := range env {
			kv := strings.SplitN(kvp, "=", 2)
			os.Setenv(kv[0], kv[1])
		}
	}(os.Environ())
	for k, v := range vars {
		if vv, ok := v.(string); ok {
			os.Setenv(k, vv)
		}
	}
	out, err := posix.ExpandEnv(p.Template)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(out))
	return err
}

func (p *Processor) executeTmpl(w io.Writer, vars map[string]interface{}) error {
	tmpl, err := template.New("").Funcs(p.funcMap).Delims(p.delimLeft, p.delimRight).Parse(p.Template)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, vars)
}
