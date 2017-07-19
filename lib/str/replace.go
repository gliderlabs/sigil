package str

import (
	"strings"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func (_ Module) Replace(old, new string, in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, sigil.ErrInputInvalid
	}
	return strings.Replace(inStr, old, new, -1), nil
}
