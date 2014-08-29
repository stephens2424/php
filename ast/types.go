package ast

import (
	"strings"
)

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
	Function

	Numeric = Float | Integer
	Unknown = String | Integer | Float | Boolean | Null | Resource | Array | Object | Function
)

var typeMap = map[Type]string{
	String:   "string",
	Integer:  "integer",
	Float:    "float",
	Boolean:  "boolean",
	Null:     "null",
	Resource: "resource",
	Array:    "array",
	Object:   "object",
	Function: "function",
}

func (t Type) Contains(typ Type) bool {
	return t&typ != 0
}

func (t Type) List() []Type {
	list := make([]Type, 0)
	for typ := range typeMap {
		if t.Contains(typ) {
			list = append(list, typ)
		}
	}
	return list
}

func (t Type) String() string {
	typeList := t.List()
	stringList := make([]string, len(typeList))
	for i, typ := range typeList {
		stringList[i] = typeMap[typ]
	}
	return strings.Join(stringList, "|")
}

type KeyType Type

const (
	StringKey  KeyType = KeyType(String)
	IntegerKey KeyType = KeyType(Integer)
)

type ArrayType struct {
	KeyType   KeyType
	ValueType Type
}

type ObjectType struct {
	Class             *Class
	DynamicProperties []Variable
}
