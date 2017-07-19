package var_

import (
	"testing"
)

func TestDefault(t *testing.T) {
	t.Parallel()
	for _, tt := range []varTest{
		{[]interface{}{"yes", nil}, "yes", nil},
		{[]interface{}{"yes", ""}, "yes", nil},
		{[]interface{}{2, ""}, 2, nil},
		{[]interface{}{false, nil}, false, nil},
		{[]interface{}{"yes", "no"}, "no", nil},
		{[]interface{}{2, 3}, 3, nil},
		{[]interface{}{false, true}, true, nil},
	} {
		got := Module{}.Default(tt.args[0], tt.args[1])
		if got != tt.expected {
			t.Errorf("Default(%#v): expected %#v, actual %#v", tt.args, tt.expected, got)
		}
	}
}
