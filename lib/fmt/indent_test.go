package fmt

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestIndent(t *testing.T) {
	for _, tt := range []fmtTest{
		{[]interface{}{"\t", "Hello"}, "\tHello", nil},
		{[]interface{}{"\t\t", "Hello\nWorld"}, "\t\tHello\n\t\tWorld", nil},
		{[]interface{}{"  ", "Hello\nWorld"}, "  Hello\n  World", nil},
		{[]interface{}{"    ", "Hello\n\nWorld"}, "    Hello\n\n    World", nil},
		{[]interface{}{"", 1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{"", true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Indent(tt.args[0].(string), tt.args[1])
		if got != tt.expected {
			t.Errorf("Indent(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Indent(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
