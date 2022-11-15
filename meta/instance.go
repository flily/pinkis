package meta

import (
	"reflect"
	"unicode"
)

// Is v a pointer
func IsPointer(v interface{}) bool {
	value := reflect.ValueOf(v)
	return value.Kind() == reflect.Ptr
}

// Is v a struct or a pointer to a struct
func IsStruct(v interface{}) bool {
	instance := InstanceOf(v)

	return instance.Kind() == reflect.Struct
}

func IsNilValue(value reflect.Value) (bool, bool) {
	switch value.Kind() {
	case reflect.Invalid:
		return true, false

	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return false, value.IsNil()

	default:
		return false, false
	}
}

func IsNil(v interface{}) (bool, bool) {
	value := reflect.ValueOf(v)
	return IsNilValue(value)
}

func IsValueInstance(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Interface:
		return false

	case reflect.Ptr:
		return value.IsNil()

	default:
		return true
	}
}

func NewPointerOf(t reflect.Type) reflect.Value {
	pointer := reflect.New(t)
	return pointer
}

func NewPointerTo(value reflect.Value) reflect.Value {
	vt := value.Type()
	addr := reflect.New(vt)
	addr.Elem().Set(value)
	return addr
}

const (
	ElemInterface = 0
	ElemPointer   = 1
)

type ValueReferenceInfo struct {
	ElemType   int
	SourceType reflect.Type
}

func ValueToInstance(value reflect.Value) (reflect.Value, []ValueReferenceInfo) {
	refChain := make([]ValueReferenceInfo, 0, 4)

	for !IsValueInstance(value) {
		info := ValueReferenceInfo{
			SourceType: value.Type(),
		}

		switch value.Kind() {
		case reflect.Ptr:
			info.ElemType = ElemPointer

		case reflect.Interface:
			info.ElemType = ElemInterface
		}

		value = value.Elem()
		refChain = append(refChain, info)
	}

	return value, refChain
}

func OriginOf(value reflect.Value, chain []ValueReferenceInfo) reflect.Value {
	if len(chain) <= 0 {
		return value
	}

	for i := len(chain) - 1; i >= 0; i-- {
		info := chain[i]

		switch info.ElemType {
		case ElemPointer:
			value = NewPointerTo(value)

		case ElemInterface:
			value = value.Convert(info.SourceType)
		}
	}

	return value
}

// Get the instance of a value, dereference all levels.
// Return final value, and dereferenced chain.
func ValueInstanceChainOf(data interface{}) (reflect.Value, []ValueReferenceInfo) {
	value := reflect.ValueOf(data)
	return ValueToInstance(value)
}

// Get the instance of a value, dereference all levels.
func ValueInstanceOf(value reflect.Value) reflect.Value {
	final, _ := ValueToInstance(value)
	return final
}

// Get the actual instance of a value, dereference all levels.
// If a nil pointer is got, return the reflect.Value represents this pointer. If an untyped nil is
// got, return an invalid reflect.Value, which is returned by reflect.ValueOf().
func InstanceOf(data interface{}) reflect.Value {
	value := reflect.ValueOf(data)
	return ValueInstanceOf(value)
}

func getFirstRune(name string) rune {
	if len(name) <= 0 {
		return 0
	}

	var first rune
	for _, c := range name {
		first = c
		break
	}

	return first
}

func IsExportedName(name string) bool {
	first := getFirstRune(name)

	// Rob Pike said, for Go 2 (can't do it before then), change the definition
	// to "lower case letters and _ are package local; all else is exported",
	// to support exported names in uncases languages like Japanese.
	//
	// see: https://github.com/golang/go/issues/5763#issuecomment-66081539

	// return !unicode.IsLower(first) && first != '_'
	return unicode.IsUpper(first)
}

func NewValueOfType(t reflect.Type) reflect.Value {
	pointer := NewPointerOf(t)
	return pointer.Elem()
}

func NewUntypedNil() reflect.Value {
	return reflect.Value{}
}

func NewTypedNil(t reflect.Type) reflect.Value {
	return reflect.Zero(t)
}

func NewValueOfValue(v reflect.Value) reflect.Value {
	isUntypedNil, isTypedNil := IsNilValue(v)
	if isUntypedNil {
		return NewUntypedNil()
	}

	if isTypedNil {
		return NewTypedNil(v.Type())
	}

	return NewValueOfType(v.Type())
}
