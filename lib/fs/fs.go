package fs

import "github.com/spf13/afero"

type Module struct {
	fs afero.Fs
}

func ModuleFunc(fs afero.Fs) func() Module {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return func() Module {
		return Module{fs}
	}
}
