package b64

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	for _, tt := range []b64Test{
		{[]interface{}{"Hello"}, "SGVsbG8=", nil},
		{[]interface{}{[]byte("Hello")}, "SGVsbG8=", nil},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Encode(tt.args[0])
		if got != tt.expected {
			t.Errorf("Encode(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Encode(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
