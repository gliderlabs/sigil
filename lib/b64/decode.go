package b64

import (
	"encoding/base64"
	"errors"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

var (
	ErrDecode = errors.New("decode error")
)

func (_ Module) Decode(in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(inStr)
	if err != nil {
		return "", ErrDecode
	}
	return string(b), nil
}
