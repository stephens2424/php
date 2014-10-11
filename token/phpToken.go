package token

import (
	"sort"
	"strconv"
)

type Token int

const (
	EOF Token = iota
	HTML
	PHP
	PHPBegin
	PHPEnd
	PHPToken
	Error
	Space
	Function
	Static
	Self
	Parent
	Final
	FunctionName
	TypeHint
	VariableOperator
	BlockBegin
	BlockEnd
	Global

	Namespace
	Use

	Comment

	IgnoreErrorOperator

	Return
	ArgumentType
	ArgumentName
	Comma
	StatementEnd
	Echo
	Print

	If
	Else
	ElseIf
	For
	Foreach
	EndIf
	EndFor
	EndForeach
	EndWhile
	EndSwitch
	AsOperator
	While
	Continue
	Break
	Do
	OpenParen
	CloseParen
	Switch
	Case
	Default

	Try
	Catch
	Finally
	Throw

	Class
	Abstract
	Private
	Public
	Protected
	Interface
	Implements
	Extends
	NewOperator
	Const

	Null
	StringLiteral
	NumberLiteral
	BooleanLiteral

	ShellCommand

	Identifier

	AssignmentOperator
	NegationOperator
	AdditionOperator
	SubtractionOperator
	MultOperator
	ConcatenationOperator
	UnaryOperator
	ComparisonOperator
	InstanceofOperator

	AndOperator
	OrOperator
	WrittenAndOperator
	WrittenXorOperator
	WrittenOrOperator

	ObjectOperator
	ScopeResolutionOperator

	CastOperator

	Var
	Array
	ArrayKeyOperator
	ArrayLookupOperatorLeft
	ArrayLookupOperatorRight
	List
	BitwiseShiftOperator
	StrongEqualityOperator
	StrongNotEqualityOperator
	EqualityOperator
	NotEqualityOperator
	AmpersandOperator
	BitwiseXorOperator
	BitwiseOrOperator
	BitwiseNotOperator
	TernaryOperator1
	TernaryOperator2

	Declare

	Include
	Exit
)

var tokens = []string{
	HTML:             "HTML",
	PHP:              "PHP",
	PHPBegin:         "PHP Begin",
	PHPEnd:           "PHP End",
	PHPToken:         "PHP Token",
	EOF:              "EOF",
	Error:            "Error",
	Space:            "Space",
	Function:         "Function",
	Static:           "static",
	Self:             "self",
	Parent:           "parent",
	Final:            "final",
	FunctionName:     "Function Name",
	TypeHint:         "Function Type Hint",
	VariableOperator: "$",
	BlockBegin:       "Block Begin",
	BlockEnd:         "Block End",

	Global:       "global",
	Return:       "Return",
	ArgumentType: "Function Argument Type",
	ArgumentName: "Function Argument Name",
	Comma:        "Function Argument Separator",
	StatementEnd: ";",
	Echo:         "echo",
	Print:        "Print",

	Namespace: "namespace",
	Use:       "use",

	IgnoreErrorOperator: "@",

	If:         "If",
	Else:       "Else",
	ElseIf:     "ElseIf",
	For:        "for",
	Foreach:    "foreach",
	Switch:     "switch",
	Case:       "case",
	Default:    "default",
	AsOperator: "as",
	While:      "while",
	Do:         "do",
	OpenParen:  "open-paren",
	CloseParen: "close-paren",
	Continue:   "continue",
	Break:      "break",
	Null:       "null",

	Comment: "/* */",

	Try:     "try",
	Catch:   "catch",
	Finally: "finally",
	Throw:   "throw",

	Class:       "Class",
	Const:       "Const",
	Abstract:    "abstract",
	Private:     "Private",
	Protected:   "Protected",
	Public:      "Public",
	Interface:   "Interface",
	Implements:  "implements",
	Extends:     "extends",
	NewOperator: "new",

	ShellCommand:   "`",
	StringLiteral:  "string-literal",
	NumberLiteral:  "number-literal",
	BooleanLiteral: "bool-literal",

	Identifier: "identifier",

	AssignmentOperator:      "=",
	NegationOperator:        "!",
	AdditionOperator:        "+",
	SubtractionOperator:     "-",
	MultOperator:            "*/%",
	ConcatenationOperator:   ".",
	UnaryOperator:           "++|--",
	ComparisonOperator:      "==<>",
	ObjectOperator:          "->",
	ScopeResolutionOperator: "::",
	InstanceofOperator:      "instanceof",

	AndOperator:        "&&",
	OrOperator:         "||",
	WrittenAndOperator: "logical-and",
	WrittenXorOperator: "logical-xor",
	WrittenOrOperator:  "logical-or",
	CastOperator:       "(type)",

	List:                     "list",
	Array:                    "array",
	ArrayKeyOperator:         "=>",
	ArrayLookupOperatorLeft:  "[",
	ArrayLookupOperatorRight: "]",
	BitwiseShiftOperator:     "<<>>",
	EqualityOperator:         "!===",
	AmpersandOperator:        "&",
	BitwiseXorOperator:       "^",
	BitwiseOrOperator:        "|",
	BitwiseNotOperator:       "~",
	TernaryOperator1:         "?",
	TernaryOperator2:         ":",

	Include: "include",
	Exit:    "exit",

	Declare: "declare",
}

var TokenList []string

func init() {
	TokenList = make([]string, len(TokenMap))
	i := 0
	for token := range TokenMap {
		TokenList[i] = token
		i += 1
	}
	sort.Sort(sort.Reverse(sort.StringSlice(TokenList)))
}

// TokenMap maps source code string tokens to  types when strings can
// be represented directly. Not all  types will be represented here.
var TokenMap = map[string]Token{
	"class":        Class,
	"clone":        UnaryOperator,
	"const":        Const,
	"abstract":     Abstract,
	"interface":    Interface,
	"implements":   Implements,
	"extends":      Extends,
	"new":          NewOperator,
	"if":           If,
	"else":         Else,
	"elseif":       ElseIf,
	"while":        While,
	"do":           Do,
	"for":          For,
	"foreach":      Foreach,
	"switch":       Switch,
	"endif;":       EndIf,
	"endif":        EndIf,
	"endfor;":      EndFor,
	"endforeach;":  EndForeach,
	"endforeach":   EndForeach,
	"endwhile;":    EndWhile,
	"endwhile":     EndWhile,
	"endswitch;":   EndSwitch,
	"endswitch":    EndSwitch,
	"case":         Case,
	"break":        Break,
	"continue":     Continue,
	"default":      Default,
	"function":     Function,
	"static":       Static,
	"final":        Final,
	"self":         Self,
	"parent":       Parent,
	"return":       Return,
	"{":            BlockBegin,
	"}":            BlockEnd,
	";":            StatementEnd,
	"(":            OpenParen,
	")":            CloseParen,
	",":            Comma,
	"echo":         Echo,
	"print":        Print,
	"throw":        Throw,
	"try":          Try,
	"catch":        Catch,
	"finally":      Finally,
	"private":      Private,
	"public":       Public,
	"protected":    Protected,
	"true":         BooleanLiteral,
	"false":        BooleanLiteral,
	"instanceof":   InstanceofOperator,
	"global":       Global,
	"list":         List,
	"array":        Array,
	"exit":         Exit,
	"include":      Include,
	"include_once": Include,
	"require":      Include,
	"require_once": Include,
	"@":            IgnoreErrorOperator,
	"null":         Null,
	"NULL":         Null,
	"var":          Var,

	"use":       Use,
	"namespace": Namespace,

	"(int)":     CastOperator,
	"(integer)": CastOperator,
	"(bool)":    CastOperator,
	"(boolean)": CastOperator,
	"(float)":   CastOperator,
	"(double)":  CastOperator,
	"(real)":    CastOperator,
	"(string)":  CastOperator,
	"(array)":   CastOperator,
	"(object)":  CastOperator,
	"(unset)":   CastOperator,

	"/*": Comment,
	"*/": Comment,
	"//": Comment,
	"#":  Comment,

	"->": ObjectOperator,
	"::": ScopeResolutionOperator,

	"+=":  AssignmentOperator,
	"-=":  AssignmentOperator,
	"*=":  AssignmentOperator,
	"/=":  AssignmentOperator,
	".=":  AssignmentOperator,
	"%=":  AssignmentOperator,
	"&=":  AssignmentOperator,
	"|=":  AssignmentOperator,
	"^=":  AssignmentOperator,
	"<<=": AssignmentOperator,
	">>=": AssignmentOperator,
	"=>":  ArrayKeyOperator,

	"===": ComparisonOperator,
	"==":  ComparisonOperator,
	"=":   AssignmentOperator,
	"!==": ComparisonOperator,
	"!=":  ComparisonOperator,
	"<>":  ComparisonOperator,
	"!":   NegationOperator,
	"++":  UnaryOperator,
	"--":  UnaryOperator,
	"+":   AdditionOperator,
	"-":   SubtractionOperator,
	"*":   MultOperator,
	"/":   MultOperator,
	">=":  ComparisonOperator,
	">":   ComparisonOperator,
	"<=":  ComparisonOperator,
	"<":   ComparisonOperator,
	"%":   MultOperator,
	".":   ConcatenationOperator,

	"&&":  AndOperator,
	"||":  OrOperator,
	"&":   AmpersandOperator,
	"^":   BitwiseXorOperator,
	"~":   BitwiseNotOperator,
	"|":   BitwiseOrOperator,
	"<<":  BitwiseShiftOperator,
	">>":  BitwiseShiftOperator,
	"?":   TernaryOperator1,
	":":   TernaryOperator2,
	"and": WrittenAndOperator,
	"xor": WrittenXorOperator,
	"or":  WrittenOrOperator,
	"as":  AsOperator,

	"[": ArrayLookupOperatorLeft,
	"]": ArrayLookupOperatorRight,

	"$":       VariableOperator,
	"declare": Declare,
}

func (i Token) String() string {
	TypeName := tokens[i]
	if len(TypeName) == 0 {
		return strconv.Itoa(int(i))
	}
	return TypeName
}
