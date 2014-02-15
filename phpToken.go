package php

import (
	"fmt"
	"strconv"
)

type item struct {
	typ itemType
	pos int
	val string
}

type location struct {
	Line, Col int64
	File      string
}

type itemType int

const (
	itemHTML itemType = iota
	itemPHP
	itemPHPBegin
	itemPHPEnd
	itemPHPToken
	itemEOF
	itemError
	itemSpace
	itemFunction
	itemFunctionName
	itemTypeHint
	itemIdentifier
	itemBlockBegin
	itemBlockEnd

	itemFunctionCallBegin
	itemFunctionCallEnd
	itemArgumentType
	itemArgumentName
	itemArgumentSeparator
	itemStatementEnd
	itemEcho

	itemIf
	itemElse
	itemElseIf
	itemFor
	itemForeach
	itemWhile
	itemDo
	itemOpenParen
	itemCloseParen

	itemTry
	itemCatch
	itemFinally

	itemClass
	itemPrivate
	itemPublic
	itemProtected
	itemInterface

	itemStringLiteral
	itemNumberLiteral
	itemBooleanLiteral

	itemNonVariableIdentifier

	itemAssignmentOperator
	itemNegationOperator
	itemAdditionOperator
	itemSubtractionOperator
	itemMultOperator
	itemConcatenationOperator
	itemUnaryOperator
	itemComparisonOperator
)

// itemTypeMap maps itemType to strings that may be used for debugging and error messages
var itemTypeMap = map[itemType]string{
	itemHTML:         "HTML",
	itemPHP:          "PHP",
	itemPHPBegin:     "PHP Begin",
	itemPHPEnd:       "PHP End",
	itemPHPToken:     "PHP Token",
	itemEOF:          "EOF",
	itemError:        "Error",
	itemSpace:        "Space",
	itemFunction:     "Function",
	itemFunctionName: "Function Name",
	itemTypeHint:     "Function Type Hint",
	itemIdentifier:   "Identifier",
	itemBlockBegin:   "Block Begin",
	itemBlockEnd:     "Block End",

	itemFunctionCallBegin: "Function Call Begin",
	itemFunctionCallEnd:   "Function Call End",
	itemArgumentType:      "Function Argument Type",
	itemArgumentName:      "Function Argument Name",
	itemArgumentSeparator: "Function Argument Separator",
	itemStatementEnd:      "Statement End",
	itemEcho:              "Echo",

	itemIf:         "If",
	itemElse:       "Else",
	itemElseIf:     "ElseIf",
	itemFor:        "for",
	itemForeach:    "foreach",
	itemWhile:      "while",
	itemDo:         "do",
	itemOpenParen:  "open-paren",
	itemCloseParen: "close-paren",

	itemTry:     "try",
	itemCatch:   "catch",
	itemFinally: "finally",

	itemClass:     "Class",
	itemPrivate:   "Private",
	itemProtected: "Protected",
	itemPublic:    "Public",
	itemInterface: "Interface",

	itemStringLiteral:  "sting-literal",
	itemNumberLiteral:  "number-literal",
	itemBooleanLiteral: "bool-literal",

	itemNonVariableIdentifier: "non-variable-identifier",

	itemAssignmentOperator:    "=",
	itemNegationOperator:      "!",
	itemAdditionOperator:      "+",
	itemSubtractionOperator:   "-",
	itemMultOperator:          "*/%",
	itemConcatenationOperator: ".",
	itemUnaryOperator:         "++|--",
	itemComparisonOperator:    "==<>",
}

// tokenMap maps source code string tokens to item types when strings can
// be represented directly. Not all item types will be represented here.
var tokenMap = map[string]itemType{
	"class":     itemClass,
	"interface": itemInterface,
	"if":        itemIf,
	"else":      itemElse,
	"while":     itemWhile,
	"for":       itemFor,
	"foreach":   itemForeach,
	"function":  itemFunction,
	"{":         itemBlockBegin,
	"}":         itemBlockEnd,
	";":         itemStatementEnd,
	"(":         itemOpenParen,
	")":         itemCloseParen,
	",":         itemArgumentSeparator,
	"echo":      itemEcho,
	"try":       itemTry,
	"catch":     itemCatch,
	"finally":   itemFinally,
	"private":   itemPrivate,
	"public":    itemPublic,
	"protected": itemProtected,
	"true":      itemBooleanLiteral,
	"false":     itemBooleanLiteral,

	"===": itemComparisonOperator,
	"==":  itemComparisonOperator,
	"=":   itemAssignmentOperator,
	"!==": itemComparisonOperator,
	"!=":  itemComparisonOperator,
	"!":   itemNegationOperator,
	"++":  itemUnaryOperator,
	"--":  itemUnaryOperator,
	"+":   itemAdditionOperator,
	"-":   itemSubtractionOperator,
	"*":   itemMultOperator,
	"/":   itemMultOperator,
	">=":  itemComparisonOperator,
	">":   itemComparisonOperator,
	"<=":  itemComparisonOperator,
	"<":   itemComparisonOperator,
	"%":   itemMultOperator,
	".":   itemConcatenationOperator,
}

func (i itemType) String() string {
	itemTypeName, ok := itemTypeMap[i]
	if !ok {
		return strconv.Itoa(int(i))
	}
	return itemTypeName
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%v:%.10q...", i.typ, i.val)
	}
	return fmt.Sprintf("%v:%q", i.typ, i.val)
}
