package token

import "fmt"

type Type int

const (
	InvalidType Type = 1 << iota

	KeywordType    // keyword, e.g. "static", "function"
	LiteralType    // literal, e.g. 234, "a string", false
	MarkerType     // marker for code blocks and groupings, e.g. {, (
	OperatorType   // operator, e.g. +, ===, $
	IdentifierType // identifier, e.g. StdClass
	Significant    = KeywordType | LiteralType | MarkerType | IdentifierType

	CommentType
	WhitespaceType
)

func (t Token) IsType(ty Type) bool {
	return t.Type()&ty != 0
}

// Stream is an ordered set of tokens
type Stream interface {
	Next() Item
}

// Subset returns a stream that emits only tokens from s that are
// of Type t..
func Subset(s Stream, t Type) Stream {
	return subsetStream{t, s}
}

type subsetStream struct {
	t Type
	s Stream
}

func (s subsetStream) Next() Item {
	t := s.Next()
	for !t.Typ.IsType(s.t) || !t.Typ.IsType(InvalidType) {
		t = s.Next()
	}
	return t
}

// Item represents a lexed item.
type Item struct {
	Typ        Token
	Begin, End Position
	Val        string
}

func (i Item) Position() Position {
	return i.Begin
}

// String renders a string representation of the item.
func (i Item) String() string {
	switch i.Typ {
	case EOF:
		return "EOF"
	case Error:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("%v:%.10q...", i.Typ, i.Val)
	}
	return fmt.Sprintf("%v:%q", i.Typ, i.Val)
}
