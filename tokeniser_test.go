package highlighter

import (
	"reflect"
	"testing"
)

func Test_Tokeniser_Next(t *testing.T) {
	tcs := []struct {
		input string
		n     int
		want  []Token
	}{
		{"lorem", 1, []Token{Word}},
		{"AND", 1, []Token{And}},
		{"  AND   ", 1, []Token{And}},
		{"  AND", 1, []Token{And}},
		{"OR", 1, []Token{Or}},
		{"\"lorem ipsum\"", 1, []Token{Phrase}},
		{"  \"lorem ipsum\"", 1, []Token{Phrase}},
		{"  \"lorem ipsum\"  ", 1, []Token{Phrase}},
		{"-lorem", 1, []Token{Minus}},
		{"NOT lorem", 1, []Token{Not}},
		{"!lorem", 1, []Token{Not}},
		{"~lorem", 1, []Token{Tilde}},
		{"lorem*", 1, []Token{Word}},
		{"lorem*", 2, []Token{Word, Wildcard}},
		{"(lorem", 1, []Token{LParen}},
		{")lorem", 1, []Token{RParen}},
		{")lorem", 2, []Token{RParen, Word}},
		{"~lorem", 1, []Token{Tilde}},
		{"~lorem", 2, []Token{Tilde, Word}},
		{"~\"lorem\"", 2, []Token{Tilde, Phrase}},
		{"lorem AND ipsum", 3, []Token{Word, And, Word}},
		{"   lorem    AND    ipsum   ", 3, []Token{Word, And, Word}},
		{"\"lorem\" AND \"ipsum\"", 3, []Token{Phrase, And, Phrase}},
		{"(\"lorem\" AND \"ipsum\") OR (lorem OR ipsum)", 11, []Token{LParen, Phrase, And, Phrase, RParen, Or, LParen, Word, Or, Word, RParen}},
		{"\"lorem AND", 2, []Token{Word, And}},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			tokeniser := NewTokeniser(tc.input)

			have := make([]Token, tc.n)
			for i := 0; i < tc.n; i++ {
				have[i] = tokeniser.Next()
			}

			if !reflect.DeepEqual(tc.want, have) {
				t.Errorf("Want: %v, Have: %v", tc.want, have)
			}
		})
	}
}

func Test_Trim(t *testing.T) {
	tcs := []struct {
		input string
		want  string
	}{
		{"lorem epsum", "lorem epsum"},
		{"   lorem epsum", "lorem epsum"},
		{"\t\tlorem epsum", "lorem epsum"},
		{"\t\nlorem epsum", "lorem epsum"},
		{"\nlorem epsum", "lorem epsum"},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			runes := []rune(tc.input)
			have := trim(runes)
			if string(have) != tc.want {
				t.Errorf("Want: %s, Have: %s", tc.want, string(have))
			}
		})
	}
}
