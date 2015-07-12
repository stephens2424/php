package token

import "fmt"

// Item represents a lexed item.
type Item struct {
	Typ        Token
	Begin, End Position
	Val        string
}

func NewItem(t Token, v string) Item {
	return Item{
		Typ: t,
		Val: v,
	}
}

func Keyword(t Token) Item {
	return Item{
		Typ: t,
		Val: fmt.Sprint(t),
	}
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
