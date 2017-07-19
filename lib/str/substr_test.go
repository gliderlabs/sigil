package str

import (
	"testing"

	"github.com/gliderlabs/sigil/pkg/sigil"
)

func TestSubstr(t *testing.T) {
	for _, tt := range []strTest{
		{[]interface{}{":", "abc"}, "abc", nil},
		{[]interface{}{":2", "abc"}, "ab", nil},
		{[]interface{}{"1:2", "abc"}, "bc", nil},
		{[]interface{}{"2:", "abcdef"}, "cdef", nil},
		{[]interface{}{"-3:2", "abcdef"}, "de", nil},
		{[]interface{}{"-3:", "abcdef"}, "def", nil},
		{[]interface{}{"", "abc"}, nil, ErrSliceIndex},
		{[]interface{}{"3", "abc"}, nil, ErrSliceIndex},
		{[]interface{}{"", 1}, nil, sigil.ErrInputInvalid},
		{[]interface{}{"", true}, nil, sigil.ErrInputInvalid},
	} {
		got, err := Module{}.Substr(tt.args[0].(string), tt.args[1])
		if got != tt.expected {
			t.Errorf("Substr(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
		if err != tt.expectedErr {
			t.Errorf("Substr(%#v): expected error %#v, actual %#v", tt.args, tt.expectedErr, err)
		}
	}
}
