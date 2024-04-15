package highlighter

import (
	"errors"
	"fmt"
)

var ErrorUnexpectedEof = errors.New("unexpected eof")
var ErrorUnexpectedToken = errors.New("unexpected token")

var stopwords = []string{
	"a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "if", "in", "into", "is",
	"it", "no", "not", "of", "on", "or", "such", "that", "the", "their", "then", "there",
	"these", "they", "this", "to", "was", "will", "with",
}

func Highlight(inputs []string, expr string) {
}

func parse(expr string) ([]string, error) {
	collect := []string{}
	tokeniser := newTokeniser(expr)

	if err := parseExpr(&tokeniser, &collect, false); err != nil {
		return collect, err
	}

	return collect, nil
}

func parseExpr(tokeniser *tokeniser, collect *[]string, parentNot bool) error {
	not := false
	for td := tokeniser.peek(); td.token != tokenEof; td = tokeniser.peek() {
		// ignore ~, *
		// if a subexpression has a leading NOT, then any leading NOTs within it should be ignored
		if td.token == tokenNot {
			not = !not
		} else if ((parentNot && not) || !parentNot && !not) &&
			(td.token == tokenWord || td.token == tokenPhrase) {
			*collect = append(*collect, td.data)
		} else if td.token == tokenAnd || td.token == tokenOr {
			// reset not at connectives
			not = false
		} else if td.token == tokenLParen {
			// if parentNot && not, then they cancel each other out
			childNot := not
			if parentNot && not {
				childNot = false
			}

			return parens(tokeniser, collect, childNot, parseExpr)
		}

		tokeniser.advance(td)
	}

	return nil
}

func parens(tokeniser *tokeniser, collect *[]string, parentNot bool, f func(*tokeniser, *[]string, bool) error) error {
	l := tokeniser.peek()
	if l.token != tokenLParen {
		return fmt.Errorf("%w: %v", ErrorUnexpectedToken, l.token)
	}
	tokeniser.advance(l)

	if err := f(tokeniser, collect, parentNot); err != nil {
		return err
	}

	r := tokeniser.peek()
	if r.token != tokenLParen {
		return fmt.Errorf("%w: %v", ErrorUnexpectedToken, r.token)
	}
	tokeniser.advance(r)

	return nil
}
