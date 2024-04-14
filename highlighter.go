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

// Accepts MATCH(column,...) against (expression)
// Returns (i, len) pairs of matches
// Optional apply (eg wrap in <b></b>, remove, replace etc)
func Highlight(inputs []string, expr string) {
	// parse the expression and find the words/phrases in the inputs
}

func parse(expr string) ([]string, error) {
	collect := []string{}
	tokeniser := newTokeniser(expr)

	if err := parseExpr(&tokeniser, &collect, false); err != nil {
		return collect, err
	}

	return collect, nil
}

// TODO: see if this can be made nicer
func parseExpr(tokeniser *tokeniser, collect *[]string, parentNot bool) error {
	// if parentNot, then any NOTs encountered in the subexpression should be accepted
	// if !parentNot, then any NOTs encountered in the subexpression should be ignored

	// TODO
	// ignore ~, *
	// need to be mindful of NOT and subexpressions
	not := false
	for tokeniser.peek().token != tokenEof {
		td := tokeniser.peek()

		// if a subexpression has a leading NOT, then any expressions within it that has a leading
		// NOT should be treated as if it didn't have that NOT
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

		tokeniser.next()
	}

	return nil
}

func parens(tokeniser *tokeniser, collect *[]string, parentNot bool, f func(*tokeniser, *[]string, bool) error) error {
	l := tokeniser.next()
	if l.token != tokenLParen {
		return fmt.Errorf("%w: %v", ErrorUnexpectedToken, l.token)
	}

	if err := f(tokeniser, collect, parentNot); err != nil {
		return err
	}

	r := tokeniser.next()
	if r.token != tokenLParen {
		return fmt.Errorf("%w: %v", ErrorUnexpectedToken, r.token)
	}

	return nil
}
