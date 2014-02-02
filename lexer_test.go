package php

import "testing"

var testFile = `<?php

session_start();

?>
<html>
<? echo something(); ?>
</html>`

func TestPHPLexer(t *testing.T) {
	l := newLexer(testFile)
	go l.run()

	item := l.nextItem()
	if item.typ != itemPHPBegin {
		t.Fatal("Did not correctly parse opening tag", item.typ)
	}

	item = l.nextItem()
	if item.typ != itemPHP {
		t.Fatal("Did not correctly parse php", item.typ)
	}
	if item.val != "\n\nsession_start();\n\n" {
		t.Fatal("Did not correctly parse php", item)
	}

	item = l.nextItem()
	if item.typ != itemPHPEnd {
		t.Fatal("Did not correctly parse ending tag", item.typ)
	}

	item = l.nextItem()
	if item.typ != itemHTML {
		t.Fatal("Did not correctly parse html", item.typ)
	}
	if item.val != "\n<html>\n" {
		t.Fatal("Did not correctly parse html", item)
	}

	item = l.nextItem()
	if item.typ != itemPHPBegin {
		t.Fatal("Did not correctly parse php begin", item.typ)
	}

	item = l.nextItem()
	if item.typ != itemPHP {
		t.Fatal("Did not correctly parse php", item.typ)
	}
	if item.val != " echo something(); " {
		t.Fatal("Did not correctly parse php", item)
	}

	item = l.nextItem()
	if item.typ != itemPHPEnd {
		t.Fatal("Did not correctly parse php end", item)
	}
	if item.val != "?>" {
		t.Fatal("Did not correctly parse php end", item)
	}

	item = l.nextItem()
	if item.typ != itemHTML {
		t.Fatal("Did not correctly parse html", item)
	}
	if item.val != "\n</html>" {
		t.Fatal("Did not correctly parse html", item)
	}

	item = l.nextItem()
	if item.typ != itemEOF {
		t.Fatal("Did not correctly parse eof", item)
	}
}
