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

	itemOperator //catchall for operators
	itemFunction
	itemFunctionName
	itemTypeHint
	itemIdentifier
	itemBlockBegin
	itemBlockEnd

	itemFunctionCallBegin
	itemFunctionCallEnd
	itemStatementEnd
	itemArgumentListBegin
	itemArgumentType
	itemArgumentName
	itemArgumentSeparator
	itemArgumentListEnd
	itemEcho
)

var itemTypeMap = map[itemType]string{
	itemHTML:         "HTML",
	itemPHP:          "PHP",
	itemPHPBegin:     "PHP Begin",
	itemPHPEnd:       "PHP End",
	itemPHPToken:     "PHP Token",
	itemEOF:          "EOF",
	itemError:        "Error",
	itemSpace:        "Space",
	itemOperator:     "Operator",
	itemFunction:     "Function",
	itemFunctionName: "Function Name",
	itemTypeHint:     "Function Type Hint",
	itemIdentifier:   "Identifier",
	itemBlockBegin:   "Block Begin",
	itemBlockEnd:     "Block End",

	itemFunctionCallBegin: "Function Call Begin",
	itemFunctionCallEnd:   "Function Call End",
	itemArgumentListBegin: "Function Argument List Begin",
	itemArgumentType:      "Function Argument Type",
	itemArgumentName:      "Function Argument Name",
	itemArgumentSeparator: "Function Argument Separator",
	itemArgumentListEnd:   "Function Argument List End",
	itemStatementEnd:      "Statement End",
	itemEcho:              "Echo",
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
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}
