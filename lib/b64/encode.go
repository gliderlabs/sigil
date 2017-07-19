package b64

import (
	"encoding/base64"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func (_ Module) Encode(in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.EncodeToString([]byte(inStr)), nil
}
