package str

import "testing"
import "github.com/gliderlabs/sigil/pkg/sigil"

func TestLower(t *testing.T) {
	for _, tt := range []strTest{
		{[]interface{}{"lowercase string"}, "lowercase string", nil},
		{[]interface{}{"LowerCase String"}, "lowercase string", nil},
		{[]interface{}{"ALLCAPS"}, "allcaps", nil},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Lower(tt.args[0])
		if got != tt.expected {
			t.Errorf("Lower(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Lower(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
