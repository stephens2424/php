package token

import (
	"fmt"
	"testing"
)

func TestTokens(t *testing.T) {
	for i := 0; i < int(maxToken); i++ {
		if len(tokens)-1 < i {
			t.Errorf("token %v has no string in tokens slice", i)
			continue
		}
		str := tokens[i]
		if str == "" {
			t.Errorf("token %v has no string in tokens slice", i)
			continue
		}
	}
}

func TestTokenTypes(t *testing.T) {
	for i := 0; i < int(maxToken); i++ {
		_, ok := tokenTypes[Token(i)]
		if !ok {
			t.Errorf(fmt.Sprintf("token %q without type", Token(i).String()))
		}
	}
}
