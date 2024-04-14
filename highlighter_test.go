package highlighter

import (
	"reflect"
	"testing"
)

func Test_Parse(t *testing.T) {
	tcs := []struct {
		input string
		want  []string
	}{
		{"", []string{}},
		{"     ", []string{}},
		{"lorem AND ipsum", []string{"lorem", "ipsum"}},
		{"   lorem    AND    ipsum   ", []string{"lorem", "ipsum"}},
		// {"!lorem AND !ipsum", []string{}},
		// {"!(lorem AND ipsum)", []string{}},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			have := parse(tc.input)

			if !reflect.DeepEqual(tc.want, have) {
				t.Errorf("Want: %v, Have: %v", tc.want, have)
			}
		})
	}
}
