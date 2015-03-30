package togo

import "go/token"

var toGoOperators = map[string]token.Token{
	"=": token.ASSIGN,

	"==":  token.EQL,
	"===": token.EQL,
	"<":   token.LSS,
	"<=":  token.LAND,
	">":   token.GTR,
	">=":  token.GEQ,

	"++": token.INC,
	"--": token.DEC,
}

func ToGoOperator(php string) token.Token {
	t, ok := toGoOperators[php]
	if !ok {
		return token.ILLEGAL
	}
	return t
}
