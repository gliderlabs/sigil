package str

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestReplace(t *testing.T) {
	for _, tt := range []strTest{
		{[]interface{}{"abc", "ABC", "abcdef"}, "ABCdef", nil},
		{[]interface{}{"b", "BB", "abcdefabc"}, "aBBcdefaBBc", nil},
		{[]interface{}{"", "", 1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{"", "", true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Replace(tt.args[0].(string), tt.args[1].(string), tt.args[2])
		if got != tt.expected {
			t.Errorf("Replace(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Replace(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
