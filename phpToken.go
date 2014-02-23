package php

import (
	"fmt"
	"sort"
	"strconv"
)

type Item struct {
	typ ItemType
	pos Location
	val string
}

type Location struct {
	Pos  int
	Line int
	File string
}

type ItemType int

const (
	itemHTML ItemType = iota
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

	itemReturn
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
	itemThrow

	itemClass
	itemPrivate
	itemPublic
	itemProtected
	itemInterface
	itemImplements
	itemExtends
	itemNewOperator

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
	itemInstanceofOperator

	itemAndOperator
	itemOrOperator
	itemWrittenAndOperator
	itemWrittenXorOperator
	itemWrittenOrOperator

	itemObjectOperator
	itemScopeResolutionOperator

	itemCastOperator

	itemArray
	itemArrayKeyOperator
	itemArrayLookupOperator
	itemBitwiseShiftOperator
	itemEqualityOperator
	itemAmpersandOperator
	itemBitwiseXorOperator
	itemBitwiseOrOperator
	itemTernaryOperator1
	itemTernaryOperator2

	itemInclude
)

// itemTypeMap maps itemType to strings that may be used for debugging and error messages
var itemTypeMap = map[ItemType]string{
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

	itemReturn:            "Return",
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
	itemThrow:   "throw",

	itemClass:       "Class",
	itemPrivate:     "Private",
	itemProtected:   "Protected",
	itemPublic:      "Public",
	itemInterface:   "Interface",
	itemImplements:  "implements",
	itemExtends:     "extends",
	itemNewOperator: "new",

	itemStringLiteral:  "sting-literal",
	itemNumberLiteral:  "number-literal",
	itemBooleanLiteral: "bool-literal",

	itemNonVariableIdentifier: "non-variable-identifier",

	itemAssignmentOperator:      "=",
	itemNegationOperator:        "!",
	itemAdditionOperator:        "+",
	itemSubtractionOperator:     "-",
	itemMultOperator:            "*/%",
	itemConcatenationOperator:   ".",
	itemUnaryOperator:           "++|--",
	itemComparisonOperator:      "==<>",
	itemObjectOperator:          "->",
	itemScopeResolutionOperator: "::",
	itemInstanceofOperator:      "instanceof",

	itemAndOperator:        "&&",
	itemOrOperator:         "||",
	itemWrittenAndOperator: "logical-and",
	itemWrittenXorOperator: "logical-xor",
	itemWrittenOrOperator:  "logical-or",
	itemCastOperator:       "(type)",

	itemArray:                "array",
	itemArrayKeyOperator:     "=>",
	itemArrayLookupOperator:  "[]",
	itemBitwiseShiftOperator: "<<>>",
	itemEqualityOperator:     "!===",
	itemAmpersandOperator:    "&",
	itemBitwiseXorOperator:   "^",
	itemBitwiseOrOperator:    "|",
	itemTernaryOperator1:     "?",
	itemTernaryOperator2:     ":",

	itemInclude: "include",
}

var tokenList []string

func init() {
	tokenList = make([]string, len(tokenMap))
	i := 0
	for token := range tokenMap {
		tokenList[i] = token
		i += 1
	}
	sort.Sort(sort.Reverse(sort.StringSlice(tokenList)))
}

// tokenMap maps source code string tokens to item types when strings can
// be represented directly. Not all item types will be represented here.
var tokenMap = map[string]ItemType{
	"class":        itemClass,
	"interface":    itemInterface,
	"implements":   itemImplements,
	"extends":      itemExtends,
	"new":          itemNewOperator,
	"if":           itemIf,
	"else":         itemElse,
	"while":        itemWhile,
	"do":           itemDo,
	"for":          itemFor,
	"foreach":      itemForeach,
	"function":     itemFunction,
	"return":       itemReturn,
	"{":            itemBlockBegin,
	"}":            itemBlockEnd,
	";":            itemStatementEnd,
	"(":            itemOpenParen,
	")":            itemCloseParen,
	",":            itemArgumentSeparator,
	"echo":         itemEcho,
	"throw":        itemThrow,
	"try":          itemTry,
	"catch":        itemCatch,
	"finally":      itemFinally,
	"private":      itemPrivate,
	"public":       itemPublic,
	"protected":    itemProtected,
	"true":         itemBooleanLiteral,
	"false":        itemBooleanLiteral,
	"instanceof":   itemInstanceofOperator,
	"array":        itemArray,
	"include":      itemInclude,
	"include_once": itemInclude,
	"require":      itemInclude,
	"require_once": itemInclude,

	"(int)":     itemCastOperator,
	"(integer)": itemCastOperator,
	"(bool)":    itemCastOperator,
	"(boolean)": itemCastOperator,
	"(float)":   itemCastOperator,
	"(double)":  itemCastOperator,
	"(real)":    itemCastOperator,
	"(string)":  itemCastOperator,
	"(array)":   itemCastOperator,
	"(object)":  itemCastOperator,
	"(unset)":   itemCastOperator,

	"->": itemObjectOperator,
	"::": itemScopeResolutionOperator,

	"+=":  itemAssignmentOperator,
	"-=":  itemAssignmentOperator,
	"*=":  itemAssignmentOperator,
	"/=":  itemAssignmentOperator,
	".=":  itemAssignmentOperator,
	"%=":  itemAssignmentOperator,
	"&=":  itemAssignmentOperator,
	"|=":  itemAssignmentOperator,
	"^=":  itemAssignmentOperator,
	"<<=": itemAssignmentOperator,
	">>=": itemAssignmentOperator,
	"=>":  itemArrayKeyOperator,

	"===": itemComparisonOperator,
	"==":  itemComparisonOperator,
	"=":   itemAssignmentOperator,
	"!==": itemComparisonOperator,
	"!=":  itemComparisonOperator,
	"<>":  itemComparisonOperator,
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

	"&&":  itemAndOperator,
	"||":  itemOrOperator,
	"&":   itemAmpersandOperator,
	"^":   itemBitwiseXorOperator,
	"|":   itemBitwiseOrOperator,
	"<<":  itemBitwiseShiftOperator,
	">>":  itemBitwiseShiftOperator,
	"?":   itemTernaryOperator1,
	":":   itemTernaryOperator2,
	"and": itemWrittenAndOperator,
	"xor": itemWrittenXorOperator,
	"or":  itemWrittenOrOperator,

	"[": itemArrayLookupOperator,
	"]": itemArrayLookupOperator,
}

func (i ItemType) String() string {
	itemTypeName, ok := itemTypeMap[i]
	if !ok {
		return strconv.Itoa(int(i))
	}
	return itemTypeName
}

func (i Item) String() string {
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
