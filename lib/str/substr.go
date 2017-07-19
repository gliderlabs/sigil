package str

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

var (
	ErrSliceIndex = errors.New("unexpected slice index")
)

func (_ Module) Substr(sliceIndex string, in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, err
	}
	indexParts := strings.Split(sliceIndex, ":")
	if len(indexParts) != 2 {
		return nil, ErrSliceIndex
	}
	start, err := strconv.Atoi(indexParts[0])
	if err != nil {
		start = 0
	}
	if start < 0 {
		start = len(inStr) + start
	}
	length, err := strconv.Atoi(indexParts[1])
	if err != nil {
		length = len(inStr) - start
	}
	return inStr[start : start+length], nil
}
