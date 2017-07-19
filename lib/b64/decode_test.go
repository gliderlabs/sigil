package b64

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestDecode(t *testing.T) {
	t.Parallel()
	for _, tt := range []b64Test{
		{[]interface{}{"SGVsbG8="}, "Hello", nil},
		{[]interface{}{[]byte("SGVsbG8=")}, "Hello", nil},
		{[]interface{}{"123"}, "", ErrDecode},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Decode(tt.args[0])
		if got != tt.expected {
			t.Errorf("Decode(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Decode(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
