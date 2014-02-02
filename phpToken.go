package php

import (
	"fmt"
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
)

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
