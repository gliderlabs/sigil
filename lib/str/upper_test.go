package str

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestUpper(t *testing.T) {
	t.Parallel()
	for _, tt := range []strTest{
		{[]interface{}{"UPPERCASE STRING"}, "UPPERCASE STRING", nil},
		{[]interface{}{"UpperCase String"}, "UPPERCASE STRING", nil},
		{[]interface{}{"lowercase"}, "LOWERCASE", nil},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Upper(tt.args[0])
		if got != tt.expected {
			t.Errorf("Upper(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Upper(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
