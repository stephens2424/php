package togo

import "go/token"

func ToGoOperator(php string) token.Token {
	switch php {
	case "==", "===":
		return token.EQL
	}
	return token.ILLEGAL
}
