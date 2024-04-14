package highlighter

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

func parse(expr string) []string {
	collect := []string{}
	tokeniser := newTokeniser(expr)

	// TODO
	// ignore ~, *
	// need to be mindful of NOT and subexpressions
	for !tokeniser.isEmpty() {
		td := tokeniser.next()

		if td.token == tokenWord || td.token == tokenPhrase {
			collect = append(collect, td.data)
		}
	}

	return collect
}
