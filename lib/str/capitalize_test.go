package str

import (
	"bytes"
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestCapitalize(t *testing.T) {
	for _, tt := range []strTest{
		{[]interface{}{"Hello"}, "Hello", nil},
		{[]interface{}{"hello"}, "Hello", nil},
		{[]interface{}{"hello three words"}, "Hello Three Words", nil},
		{[]interface{}{"hello \nnewline"}, "Hello \nNewline", nil},
		{[]interface{}{[]byte("Hello")}, "Hello", nil},
		{[]interface{}{bytes.NewBufferString("Hello")}, "Hello", nil},
		{[]interface{}{1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Capitalize(tt.args[0])
		if got != tt.expected {
			t.Errorf("Capitalize(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Capitalize(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
