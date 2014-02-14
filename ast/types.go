package ast

type Type int

const (
	String Type = 1 << iota
	Integer
	Float
	Boolean
	Null
	Resource
	Array
	Object
)

type KeyType Type

const (
	StringKey  KeyType = iota
	IntegerKey KeyType = iota
)

type ArrayType struct {
	KeyType   KeyType
	ValueType Type
}

type ObjectType struct {
	Class             *Class
	DynamicProperties []Identifier
}
