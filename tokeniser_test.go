package highlighter

import (
	"reflect"
	"testing"
)

func Test_Tokeniser_Next(t *testing.T) {
	tcs := []struct {
		input string
		n     int
		want  []tokenData
	}{
		{"", 1, []tokenData{{token: tokenEof}}},
		{"    ", 1, []tokenData{{token: tokenEof}}},
		{"lorem", 1, []tokenData{{tokenWord, "lorem"}}},
		{"AND", 1, []tokenData{{token: tokenAnd}}},
		{"  AND   ", 1, []tokenData{{token: tokenAnd}}},
		{"  AND", 1, []tokenData{{token: tokenAnd}}},
		{"OR", 1, []tokenData{{token: tokenOr}}},
		{"lorem ipsum", 2, []tokenData{{tokenWord, "lorem"}, {tokenWord, "ipsum"}}},
		{"\"lorem ipsum\"", 1, []tokenData{{tokenPhrase, "lorem ipsum"}}},
		{"  \"lorem ipsum\"", 1, []tokenData{{tokenPhrase, "lorem ipsum"}}},
		{"  \"lorem ipsum\"  ", 1, []tokenData{{tokenPhrase, "lorem ipsum"}}},
		{"-lorem", 1, []tokenData{{token: tokenNot}}},
		{"NOT lorem", 1, []tokenData{{token: tokenNot}}},
		{"!lorem", 1, []tokenData{{token: tokenNot}}},
		{"~lorem", 1, []tokenData{{token: tokenTilde}}},
		{"lorem*", 1, []tokenData{{tokenWord, "lorem"}}},
		{"lorem*", 2, []tokenData{{tokenWord, "lorem"}, {token: tokenWildcard}}},
		{"(lorem", 1, []tokenData{{token: tokenLParen}}},
		{")lorem", 1, []tokenData{{token: tokenRParen}}},
		{")lorem", 2, []tokenData{{token: tokenRParen}, {tokenWord, "lorem"}}},
		{"~lorem", 1, []tokenData{{token: tokenTilde}}},
		{"~lorem", 2, []tokenData{{token: tokenTilde}, {tokenWord, "lorem"}}},
		{"~\"lorem\"", 2, []tokenData{{token: tokenTilde}, {tokenPhrase, "lorem"}}},
		{"!\"ipsum\"", 2, []tokenData{{token: tokenNot}, {tokenPhrase, "ipsum"}}},
		{"lorem AND ipsum", 3, []tokenData{{tokenWord, "lorem"}, {token: tokenAnd}, {tokenWord, "ipsum"}}},
		{"   lorem    AND    ipsum   ", 3, []tokenData{{tokenWord, "lorem"}, {token: tokenAnd}, {tokenWord, "ipsum"}}},
		{"\"lorem\" AND \"ipsum\"", 3, []tokenData{{tokenPhrase, "lorem"}, {token: tokenAnd}, {tokenPhrase, "ipsum"}}},
		{"(\"lorem\" AND \"ipsum\") OR (lorem OR ipsum)", 11, []tokenData{{token: tokenLParen}, {tokenPhrase, "lorem"}, {token: tokenAnd}, {tokenPhrase, "ipsum"}, {token: tokenRParen}, {token: tokenOr}, {token: tokenLParen}, {tokenWord, "lorem"}, {token: tokenOr}, {tokenWord, "ipsum"}, {token: tokenRParen}}},
		{"\"lorem AND", 2, []tokenData{{tokenWord, "lorem"}, {token: tokenAnd}}},
		{"+\"ipsum\"", 2, []tokenData{{token: tokenPlus}, {tokenPhrase, "ipsum"}}},
		{"NOT \"ipsum\"", 2, []tokenData{{token: tokenNot}, {tokenPhrase, "ipsum"}}},
		{"-\"ipsum\"", 2, []tokenData{{token: tokenNot}, {tokenPhrase, "ipsum"}}},
		{"!\"ipsum\"", 2, []tokenData{{token: tokenNot}, {tokenPhrase, "ipsum"}}},
		{"\"lorem\" AND !\"ipsum\"", 4, []tokenData{{tokenPhrase, "lorem"}, {token: tokenAnd}, {token: tokenNot}, {tokenPhrase, "ipsum"}}},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			tokeniser := newTokeniser(tc.input)

			have := make([]tokenData, tc.n)
			for i := 0; i < tc.n; i++ {
				have[i] = tokeniser.peek()
				tokeniser.advance(have[i])
			}

			if !reflect.DeepEqual(tc.want, have) {
				t.Errorf("Want: %v, Have: %v", tc.want, have)
			}
		})
	}
}

func Test_SkipWhitespace(t *testing.T) {
	tcs := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"    ", ""},
		{"lorem epsum", "lorem epsum"},
		{"   lorem epsum", "lorem epsum"},
		{"\t\tlorem epsum", "lorem epsum"},
		{"\t\nlorem epsum", "lorem epsum"},
		{"\nlorem epsum", "lorem epsum"},
		{"*", "*"},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			runes := []rune(tc.input)
			have := skipWhitespace(runes)
			if string(have) != tc.want {
				t.Errorf("Want: %s, Have: %s", tc.want, string(have))
			}
		})
	}
}
