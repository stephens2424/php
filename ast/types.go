package ast

type keyType BasicType

const (
	MixedKey   = String | Integer
	StringKey  = String
	IntegerKey = Integer
)

type ArrayType struct {
	KeyType   keyType
	ValueType Type
}

type BasicType int

const (
	Invalid BasicType = iota
	String
	Integer
	Float
	Boolean
	Null
	Resource
	Array
	Object
	Function
)

var Numeric = compoundType{Integer: struct{}{}, Float: struct{}{}}

var typeMap = map[BasicType]string{
	String:   "string",
	Integer:  "integer",
	Float:    "float",
	Boolean:  "boolean",
	Null:     "null",
	Resource: "resource",
	Array:    "array",
	Object:   "object",
	Function: "function",
	Invalid:  "invalid-type",
}

func (t BasicType) Contains(typ Type) bool {
	if bt, ok := typ.(BasicType); ok {
		return t&bt != 0
	}
	return false
}

func (t BasicType) String() string {
	if st, ok := typeMap[t]; ok {
		return st
	}
	return typeMap[Invalid]
}

func (t BasicType) Basic() []BasicType {
	return []BasicType{t}
}

func (t BasicType) Equals(o Type) bool {
	ot, ok := o.(BasicType)
	if !ok {
		return false
	}
	return ot == t
}

func (t BasicType) Single() bool {
	return t != 0 && t&(t-1) == 0
}

func (t BasicType) Union(o Type) Type {
	return compoundType{t: struct{}{}, o: struct{}{}}
}

type compoundType map[Type]struct{}

func (c compoundType) Equals(t Type) bool {
	if ct, ok := t.(compoundType); ok {
		if len(ct) != len(c) {
			return false
		}
		for it := range c {
			if _, ok := ct[it]; !ok {
				return false
			}
		}
		return true
	}
	if len(c) == 1 {
		for it := range c {
			return it.Equals(t)
		}
	}
	return false
}

func (c compoundType) Contains(t Type) bool {
	if ct, ok := t.(compoundType); ok {
		if len(ct) > len(c) {
			return false
		}
		for it := range ct {
			if _, ok := c[it]; !ok {
				return false
			}
		}
		return true
	}

	for it := range c {
		if it.Contains(t) {
			return true
		}
	}

	return false
}

// Union returns a new type that includes both the receiver and the argument.
func (c compoundType) Union(t Type) Type {
	c[t] = struct{}{}
	return c
}

// Single returns true if the receiver expresses one type and only one type.
func (c compoundType) Single() bool {
	if len(c) != 1 {
		return false
	}
	for t := range c {
		return t.Single()
	}
	return false
}

// String returns the receiver expressed as a string.
func (c compoundType) String() string {
	return ""
}

// Basic returns the basic type a type expresses.
func (c compoundType) Basic() []BasicType {
	return nil
}

type ObjectType struct {
	Class string
}

type Type interface {
	// Equals returns true if the receiver is of the same type as the argument.
	Equals(Type) bool

	// Contains returns true if the receiver contains the argument type.
	Contains(Type) bool

	// Union returns a new type that includes both the receiver and the argument.
	Union(Type) Type

	// Single returns true if the receiver expresses one type and only one type.
	Single() bool

	// String returns the receiver expressed as a string.
	String() string

	// Basic returns the basic type a type expresses.
	Basic() []BasicType
}

var Unknown = new(unknownType)

type unknownType struct{}

func (_ unknownType) Equals(t Type) bool {
	return t == Unknown
}

func (_ unknownType) Contains(t Type) bool {
	return t == Unknown
}

func (_ unknownType) Union(t Type) Type {
	return t
}

func (_ unknownType) Single() bool {
	return false
}

func (_ unknownType) String() string {
	return "unknown"
}

func (_ unknownType) Basic() []BasicType {
	return nil
}
