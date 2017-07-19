package fs

import "github.com/spf13/afero"

func (m Module) Read(filepath string) (string, error) {
	b, err := afero.ReadFile(m.fs, filepath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
