package token

import (
	"testing"
)

func TestNewList(t *testing.T) {
	list := NewList(Item{Typ: Echo}, Item{Typ: Echo})

	if len(list.Items) != 2 {
		t.Fatalf("Items should contain two items")
	}
}

func TestNext(t *testing.T) {
	list := NewList(Item{Typ: Echo}, Item{Typ: Echo})
	items := make([]Item, 0)

	for item := list.Next(); item.Typ != EOF; item = list.Next() {
		items = append(items, item)
	}

	if len(items) != 2 {
		t.Fatalf("Items should contain two items")
	}
}

func TestPeek(t *testing.T) {
	list := NewList(Item{Typ: Echo}, Item{Typ: HTML})
	list.Next()

	if list.Peek().Typ != HTML {
		t.Fatalf("Item should be type HTML")
	}
}

func TestPush(t *testing.T) {
	list := NewList()
	list.Push(Item{Typ: Echo})

	if len(list.Items) != 1 {
		t.Fatalf("Items should contain one items")
	}
}

func TestSeek(t *testing.T) {
	list := NewList(Item{Typ: Echo}, Item{Typ: HTML})

	list.Seek(1)
	if list.Next().Typ != HTML {
		t.Fatalf("Item should be type HTML")
	}

	list.Seek(0)
	if list.Next().Typ != Echo {
		t.Fatalf("Item should be type Echo")
	}
}

func TestPushStream(t *testing.T) {
	list := NewList(Item{Typ: Echo}, Item{Typ: Echo})
	list.PushStream(NewList(Item{Typ: Echo}, Item{Typ: Echo}))

	if len(list.Items) != 4 {
		t.Fatalf("Items should contain four items")
	}
}
