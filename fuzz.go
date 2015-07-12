package php

import (
	"bytes"
	"os/exec"
)

// +build gofuzz

func Fuzz(data []byte) int {
	interestingness := 0
	isValid := validPHP(data)

	if isValid {
		interestingness++
	}

	p := NewParser()
	p.Debug = true

	_, err := p.Parse("test.php", string(data))
	if err != nil {
		interestingness++
	}

	if err == nil && !isValid {
		panic("false negative: the input was invalid and not flagged")
	} else if err != nil && isValid {
		panic("false positive: the input was valid but flagged invalid")
	}

	return interestingness
}

func validPHP(data []byte) bool {
	cmd := exec.Command("php", "-l")
	cmd.Stdin = bytes.NewReader(data)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
