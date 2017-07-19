package str

import (
	"strings"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func (_ Module) Lower(in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, err
	}
	return strings.ToLower(inStr), nil
}
