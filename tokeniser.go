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

type tokenLiteral string

const (
	tokenLiteralPlus     tokenLiteral = "+"
	tokenLiteralMinus    tokenLiteral = "-"
	tokenLiteralNot      tokenLiteral = "NOT"
	tokenLiteralAnd      tokenLiteral = "AND"
	tokenLiteralOr       tokenLiteral = "OR"
	tokenLiteralLParen   tokenLiteral = "("
	tokenLiteralRParen   tokenLiteral = ")"
	tokenLiteralWildcard tokenLiteral = "*"
	tokenLiteralTilde    tokenLiteral = "~"

	tokenLiteralNotSymbol tokenLiteral = "!"
)

func (t tokenLiteral) String() string {
	return string(t)
}

var symbols = []tokenLiteral{tokenLiteralPlus, tokenLiteralMinus, tokenLiteralLParen, tokenLiteralRParen, tokenLiteralWildcard, tokenLiteralTilde, tokenLiteralNotSymbol}

type token int

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
		return tokenLiteralPlus.String()
	case tokenMinus:
		return tokenLiteralMinus.String()
	case tokenNot:
		return tokenLiteralNot.String()
	case tokenAnd:
		return tokenLiteralAnd.String()
	case tokenOr:
		return tokenLiteralOr.String()
	case tokenLParen:
		return tokenLiteralLParen.String()
	case tokenRParen:
		return tokenLiteralRParen.String()
	case tokenWildcard:
		return tokenLiteralWildcard.String()
	case tokenTilde:
		return tokenLiteralTilde.String()
	}

	panic("unreachable")
}

func newToken(s string) tokenData {
	switch s {
	case tokenLiteralPlus.String():
		return tokenData{token: tokenPlus}
	case tokenLiteralMinus.String():
		return tokenData{token: tokenMinus}
	case tokenLiteralNot.String(), tokenLiteralNotSymbol.String():
		return tokenData{token: tokenNot}
	case tokenLiteralAnd.String():
		return tokenData{token: tokenAnd}
	case tokenLiteralOr.String():
		return tokenData{token: tokenOr}
	case tokenLiteralLParen.String():
		return tokenData{token: tokenLParen}
	case tokenLiteralRParen.String():
		return tokenData{token: tokenRParen}
	case tokenLiteralWildcard.String():
		return tokenData{token: tokenWildcard}
	case tokenLiteralTilde.String():
		return tokenData{token: tokenTilde}
	}

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return tokenData{tokenPhrase, s}
	}

	return tokenData{tokenWord, s}
}

type tokenData struct {
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

func (t *tokeniser) next() tokenData {
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
		return tokenData{tokenPhrase, phrase}
	}

token:
	for i, r := range t.expr {
		if slices.Contains(symbols, tokenLiteral(r)) {
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
