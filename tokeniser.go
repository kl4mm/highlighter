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

var symbols = []string{Plus.String(), Minus.String(), LParen.String(), RParen.String(), Wildcard.String(), Tilde.String(), NotSymbol}

type Token int

const NotSymbol = "!"
const (
	Word Token = iota
	Phrase
	Plus
	Minus
	Not
	And
	Or
	LParen
	RParen
	Wildcard
	Tilde
)

func (t Token) String() string {
	switch t {
	case Word:
		return "<word>"
	case Phrase:
		return "<phrase>"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Not:
		return "NOT"
	case And:
		return "AND"
	case Or:
		return "OR"
	case LParen:
		return "("
	case RParen:
		return ")"
	case Wildcard:
		return "*"
	case Tilde:
		return "~"
	}

	panic("unreachable")
}

func NewToken(s string) Token {
	switch s {
	case Plus.String():
		return Plus
	case Minus.String():
		return Minus
	case Not.String(), NotSymbol:
		return Not
	case And.String():
		return And
	case Or.String():
		return Or
	case LParen.String():
		return LParen
	case RParen.String():
		return RParen
	case Wildcard.String():
		return Wildcard
	case Tilde.String():
		return Tilde
	}

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return Phrase
	}

	return Word
}

type Tokeniser struct {
	expr []rune
}

func NewTokeniser(expr string) Tokeniser {
	return Tokeniser{expr: []rune(expr)}
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

func (t *Tokeniser) Next() Token {
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
			goto token
		}

		// Unquote
		_ = string(t.expr[1 : to+1])
		// Quote
		// phrase := string(t.expr[:to+2])

		t.expr = t.expr[to+2:]
		return Phrase
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
	return NewToken(token)
}
