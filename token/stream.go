package token

// Stream is an ordered set of tokens
type Stream interface {
	// Next consumes and returns the next item in the stream. If there is no next
	// item, the zero value is returned.
	Next() Item

	// Previous consumes and returns the previous item in the stream. i.e.
	//
	// 	a := s.Next()
	//	b := s.Previous()
	//	c := s.Next()
	//	a == b == c // true
	//
	// If there is no previous item, the zero value is returned.
	Previous() Item
}

// List represents an ordered set of tokens.
type itemList struct {
	// Items contains all the items in the list.
	Items []Item

	// Position is the current position the set is at in the token slice.
	Position int
}

// NewList initializes a new ItemList
func NewList(t ...Item) *itemList {
	return &itemList{t, 0}
}

// Next consumes and returns the next item in the list.
func (s *itemList) Next() Item {
	if s.Position == len(s.Items) {
		return Item{}
	}

	item := s.Items[s.Position]
	s.Position++

	return item
}

func (s *itemList) Previous() Item {
	if s.Position == 0 {
		return Item{}
	}

	s.Position--
	return s.Items[s.Position]
}

func (s *itemList) Peek() Item {
	return s.Items[s.Position]
}

func (s *itemList) Push(i ...Item) {
	s.Items = append(s.Items, i...)
}

func (s *itemList) PushKeyword(t Token) {
	s.Items = append(s.Items, Keyword(t))
}

func (s *itemList) PushStream(i Stream) {
	for item := i.Next(); item.Typ != EOF; item = i.Next() {
		s.Push(item)
	}
}

func (s *itemList) Seek(position int) {
	s.Position = position
}

// Subset returns a stream that emits only tokens from s that are
// of Type t. If s is already a subset, it is re-expanded before
// the new subset is created.
func Subset(s Stream, t Type) Stream {
	if sb, ok := s.(subsetStream); ok {
		return subsetStream{t: t, s: sb.s}
	}
	return subsetStream{t: t, s: s}
}

type subsetStream struct {
	t Type
	s Stream
}

func (s subsetStream) Next() Item {
	t := s.s.Next()
	for !t.Typ.Type().Is(s.t) && !t.Typ.Type().Is(InvalidType) {
		t = s.s.Next()
	}
	return t
}

func (s subsetStream) Previous() Item {
	t := s.s.Previous()
	for !t.Typ.Type().Is(s.t) && !t.Typ.Type().Is(InvalidType) {
		t = s.s.Previous()
	}
	return t
}
