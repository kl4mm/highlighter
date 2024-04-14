package highlighter

import (
	"slices"
	"strings"
	"unicode"
)

// Accepts MATCH(column) against (expression)
// Returns (i, len) pairs of matches
// Optional apply (eg wrap in <b></b>, remove, replace etc)

var stopwords = []string{
	"a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "if", "in", "into", "is",
	"it", "no", "not", "of", "on", "or", "such", "that", "the", "their", "then", "there",
	"these", "they", "this", "to", "was", "will", "with",
}

var symbols = []string{tokenPlus.String(), tokenMinus.String(), tokenLParen.String(), tokenRParen.String(), tokenWildcard.String(), tokenTilde.String(), notSymbol}

type token int

const notSymbol = "!"
const (
	tokenWord token = iota
	tokenPhrase
	tokenPlus
	tokenMinus
	tokenNot
	tokenAnd
	tokenOr
	tokenLParen
	tokenRParen
	tokenWildcard
	tokenTilde
)

func (t token) String() string {
	switch t {
	case tokenWord:
		return "<word>"
	case tokenPhrase:
		return "<phrase>"
	case tokenPlus:
		return "+"
	case tokenMinus:
		return "-"
	case tokenNot:
		return "NOT"
	case tokenAnd:
		return "AND"
	case tokenOr:
		return "OR"
	case tokenLParen:
		return "("
	case tokenRParen:
		return ")"
	case tokenWildcard:
		return "*"
	case tokenTilde:
		return "~"
	}

	panic("unreachable")
}

func newToken(s string) TokenData {
	switch s {
	case tokenPlus.String():
		return TokenData{token: tokenPlus}
	case tokenMinus.String():
		return TokenData{token: tokenMinus}
	case tokenNot.String(), notSymbol:
		return TokenData{token: tokenNot}
	case tokenAnd.String():
		return TokenData{token: tokenAnd}
	case tokenOr.String():
		return TokenData{token: tokenOr}
	case tokenLParen.String():
		return TokenData{token: tokenLParen}
	case tokenRParen.String():
		return TokenData{token: tokenRParen}
	case tokenWildcard.String():
		return TokenData{token: tokenWildcard}
	case tokenTilde.String():
		return TokenData{token: tokenTilde}
	}

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return TokenData{tokenPhrase, s}
	}

	return TokenData{tokenWord, s}
}

type TokenData struct {
	token token
	data  string
}

type tokeniser struct {
	expr []rune
}

func newTokeniser(expr string) tokeniser {
	return tokeniser{expr: []rune(expr)}
}

func trim(runes []rune) []rune {
	from := 0
	for i, r := range runes {
		from = i
		if !unicode.IsSpace(r) {
			break
		}
	}

	return runes[from:]
}

func (t *tokeniser) next() TokenData {
	t.expr = trim(t.expr)

	to := 0
	if t.expr[0] == '"' {
		for i, r := range t.expr[1:] {
			if r == '"' {
				to = i
				break
			}
		}

		// If we couldn't find the closing quote, then ignore it and find the next token
		if to == 0 {
			t.expr = t.expr[1:]
			goto token
		}

		// Unquote
		phrase := string(t.expr[1 : to+1])
		// Quote
		// phrase := string(t.expr[:to+2])

		t.expr = t.expr[to+2:]
		return TokenData{tokenPhrase, phrase}
	}

token:
	for i, r := range t.expr {
		if slices.Contains(symbols, string(r)) {
			break
		}

		if unicode.IsSpace(r) {
			break
		}

		to = i
	}

	token := string(t.expr[:to+1])
	t.expr = t.expr[to+1:]
	return newToken(token)
}
