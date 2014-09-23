package token

type Position struct {
	Line, Column int // The position relative to other characters in the file
	Position     int // The position in bytes in the file
	File         string
}
