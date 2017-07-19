package fs

import (
	"testing"

	"github.com/spf13/afero"
)

func TestRead(t *testing.T) {
	t.Parallel()
	fs := afero.NewMemMapFs()
	fs.MkdirAll("dir", 0755)
	afero.WriteFile(fs, "file1.txt", []byte("foo"), 0644)
	afero.WriteFile(fs, "file2.txt", []byte("bar"), 0644)
	afero.WriteFile(fs, "dir/file3.txt", []byte("baz"), 0644)
	for _, tt := range []fsTest{
		{[]interface{}{"file1.txt"}, "foo", nil},
		{[]interface{}{"file2.txt"}, "bar", nil},
		{[]interface{}{"dir/file3.txt"}, "baz", nil},
		{[]interface{}{"nofile.txt"}, "", "open nofile.txt: file does not exist"},
	} {
		got, err := Module{fs}.Read(tt.args[0].(string))
		if got != tt.expected {
			t.Errorf("Read(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if (err != nil && tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.(string)) {
			t.Errorf("Read(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err.Error())
		}
	}
}
