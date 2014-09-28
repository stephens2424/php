package token

func (t Token) Type() Type {
	typ, ok := tokenTypes[t]
	if !ok {
		panic("token without type")
	}
	return typ
}

var tokenTypes = map[Token]Type{
	HTML:     LiteralType,
	PHPBegin: KeywordType,
	PHPEnd:   KeywordType,
	PHPToken: KeywordType,

	EOF:   InvalidType,
	Error: InvalidType,

	Space: WhitespaceType,

	Function:  KeywordType,
	Static:    KeywordType,
	Self:      KeywordType,
	Parent:    KeywordType,
	Final:     KeywordType,
	Global:    KeywordType,
	Return:    KeywordType,
	Namespace: KeywordType,
	Use:       KeywordType,
	Echo:      KeywordType,
	Print:     KeywordType,

	FunctionName:     IdentifierType,
	TypeHint:         IdentifierType,
	VariableOperator: OperatorType,

	Comma:        MarkerType,
	StatementEnd: MarkerType,

	BlockBegin: MarkerType,
	BlockEnd:   MarkerType,

	IgnoreErrorOperator: OperatorType,

	If:         KeywordType,
	Else:       KeywordType,
	ElseIf:     KeywordType,
	For:        KeywordType,
	Foreach:    KeywordType,
	Switch:     KeywordType,
	Case:       KeywordType,
	Default:    KeywordType,
	AsOperator: KeywordType,
	While:      KeywordType,
	Do:         KeywordType,
	Continue:   KeywordType,
	Break:      KeywordType,
	Try:        KeywordType,
	Catch:      KeywordType,
	Finally:    KeywordType,
	Throw:      KeywordType,

	OpenParen:  MarkerType,
	CloseParen: MarkerType,

	Null:    IdentifierType,
	Comment: CommentType,

	Class:       KeywordType,
	Const:       KeywordType,
	Abstract:    KeywordType,
	Private:     KeywordType,
	Protected:   KeywordType,
	Public:      KeywordType,
	Interface:   KeywordType,
	Implements:  KeywordType,
	Extends:     KeywordType,
	NewOperator: KeywordType,

	ShellCommand:   LiteralType,
	StringLiteral:  LiteralType,
	NumberLiteral:  LiteralType,
	BooleanLiteral: LiteralType,

	Identifier: IdentifierType,

	AssignmentOperator:      OperatorType,
	NegationOperator:        OperatorType,
	AdditionOperator:        OperatorType,
	SubtractionOperator:     OperatorType,
	MultOperator:            OperatorType,
	ConcatenationOperator:   OperatorType,
	UnaryOperator:           OperatorType,
	ComparisonOperator:      OperatorType,
	ObjectOperator:          OperatorType,
	ScopeResolutionOperator: OperatorType,
	InstanceofOperator:      OperatorType,
	AndOperator:             OperatorType,
	OrOperator:              OperatorType,
	WrittenAndOperator:      OperatorType,
	WrittenXorOperator:      OperatorType,
	WrittenOrOperator:       OperatorType,
	CastOperator:            OperatorType,

	List:                     KeywordType,
	Array:                    KeywordType,
	ArrayKeyOperator:         KeywordType,
	ArrayLookupOperatorLeft:  MarkerType,
	ArrayLookupOperatorRight: MarkerType,

	BitwiseShiftOperator: OperatorType,
	EqualityOperator:     OperatorType,
	AmpersandOperator:    OperatorType,
	BitwiseXorOperator:   OperatorType,
	BitwiseOrOperator:    OperatorType,
	BitwiseNotOperator:   OperatorType,
	TernaryOperator1:     OperatorType,
	TernaryOperator2:     OperatorType,

	Include: KeywordType,
	Exit:    KeywordType,

	Declare: KeywordType,
}
