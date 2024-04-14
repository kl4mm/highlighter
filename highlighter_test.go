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
		{"!lorem AND !ipsum", []string{}},
		{"!(lorem AND ipsum)", []string{}},
		{"!(!lorem AND !ipsum)", []string{"lorem", "ipsum"}},
		{"!(!(lorem AND ipsum))", []string{"lorem", "ipsum"}},
		{"!(!(!lorem AND !ipsum))", []string{}},
		{"!lorem AND ipsum", []string{"ipsum"}},
		{"!lorem AND !(ipsum AND !lorem)", []string{"lorem"}},
		{"!lorem AND !(ipsum AND !(dolor sit amet))", []string{"dolor", "sit", "amet"}},
		{"!lorem AND !(ipsum AND !(dolor AND !(sit amet)))", []string{"dolor"}},
		{"!lorem AND (ipsum AND !(dolor AND !(sit amet)))", []string{"ipsum", "sit", "amet"}},
		{"!lorem AND (ipsum AND !!(dolor AND !(sit amet)))", []string{"ipsum", "dolor"}},
		{"!\"lorem\" AND (\"ipsum\" AND !!(\"dolor\" AND !(\"sit\" \"amet\")))", []string{"ipsum", "dolor"}},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			have, _ := parse(tc.input)

			if !reflect.DeepEqual(tc.want, have) {
				t.Errorf("Want: %v, Have: %v", tc.want, have)
			}
		})
	}
}
