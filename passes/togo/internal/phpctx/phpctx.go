package phpctx

import (
	"errors"
	"io"
	"reflect"
)

var (
	ErrNoStruct     = errors.New("receiver is not a struct")
	ErrMissingField = errors.New("struct is missing field")
	ErrNotSet       = errors.New("dynamic value is not set")
)

var zero = reflect.Value{}

type PHPContext struct {
	Echo          io.Writer
	dynamicValues map[string]*interface{}
}

func (ctx PHPContext) SetDynamic(name string, value interface{}) {
	ctx.dynamicValues[name] = &value
}

func (ctx PHPContext) GetDynamic(name string) (interface{}, error) {
	v, ok := ctx.dynamicValues[name]
	if !ok {
		return nil, ErrNotSet
	}

	if v == nil {
		return nil, nil
	}

	return *v, nil
}

func GetDynamicProperty(rcvr interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(rcvr)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(field)
		if f == zero {
			return nil, ErrMissingField
		}
		return f.Interface(), nil
	}
	return nil, ErrNoStruct
}
