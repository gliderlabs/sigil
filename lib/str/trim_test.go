package str

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestTrim(t *testing.T) {
	t.Parallel()
	for _, tt := range []strTest{
		{[]interface{}{"trimmed text"}, "trimmed text", nil},
		{[]interface{}{"  untrimmed text "}, "untrimmed text", nil},
		{[]interface{}{"\nnewline \ntrimmed\n"}, "newline \ntrimmed", nil},
		{[]interface{}{"newline \nspaces \n  \n"}, "newline \nspaces", nil},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Trim(tt.args[0])
		if got != tt.expected {
			t.Errorf("Trim(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Trim(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
