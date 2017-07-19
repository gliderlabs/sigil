package sigil

import (
	"errors"
	"fmt"
)

var (
	ErrInputInvalid = errors.New("input argument invalid")
)

func InputString(in interface{}) (string, error) {
	switch obj := in.(type) {
	case string:
		return obj, nil
	case []byte:
		return string(obj), nil
	case fmt.Stringer:
		return obj.String(), nil
	default:
		return "", ErrInputInvalid
	}
}
