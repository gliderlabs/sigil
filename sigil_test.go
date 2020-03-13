package sigil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLookPath(t *testing.T) {
	for _, tmpDir := range []string{"", "."} {
		for _, relative := range []bool{true, false} {
			for _, exists := range []bool{true, false} {
				// setup
				f, err := ioutil.TempFile(tmpDir, "sigiltest")
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
