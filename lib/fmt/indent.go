package fmt

import (
	"strings"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func (_ Module) Indent(indent string, in interface{}) (interface{}, error) {
	inStr, err := sigil.InputString(in)
	if err != nil {
		return nil, err
	}
	var indented []string
	lines := strings.Split(inStr, "\n")
	indented = append(indented, indent+lines[0])
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			if line != "" {
				indented = append(indented, indent+line)
			} else {
				indented = append(indented, line)
			}
		}
	}
	return strings.Join(indented, "\n"), nil
}
