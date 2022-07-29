package meta

import (
	"reflect"
	"unsafe"
)

// Make a new unsafe reference to the value.
// Use to copy value included unexported field
func MakeUnsafeRef(data reflect.Value) reflect.Value {
	var result reflect.Value

	switch data.Kind() {
	case reflect.Ptr:
		instance := data.Elem()
		instancePointer := unsafe.Pointer(instance.UnsafeAddr())
		unsafePointer := reflect.NewAt(instance.Type(), instancePointer)
		result = unsafePointer

	case reflect.Array, reflect.Slice, reflect.Map:
		result = data

	default:
		instance := data
		instancePointer := unsafe.Pointer(instance.UnsafeAddr())
		unsafePointer := reflect.NewAt(instance.Type(), instancePointer)
		unsafeInstance := unsafePointer.Elem()
		result = unsafeInstance
	}

	return result
}

func UnsafeValueSet(target reflect.Value, source reflect.Value) {
	targetPointer := unsafe.Pointer(target.UnsafeAddr())
	targetUnsafe := reflect.NewAt(target.Type(), targetPointer).Elem()

	unsafeSource := MakeUnsafeRef(source)
	targetUnsafe.Set(unsafeSource)
}

func duplicateValueForArray(data reflect.Value) (reflect.Value, error) {
	newSlice := reflect.MakeSlice(data.Type(), data.Len(), data.Cap())
	for i := 0; i < data.Len(); i++ {
		itemValue := data.Index(i)
		itemCopy, err := duplicateValueInstance(itemValue)
		if err != nil {
			return reflect.Value{}, err
		}

		item := newSlice.Index(i)
		UnsafeValueSet(item, itemCopy)
	}

	return newSlice, nil
}

func duplicateValueForMap(data reflect.Value) (reflect.Value, error) {
	newMap := reflect.MakeMapWithSize(data.Type(), data.Len())
	iter := data.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyCopy, err := duplicateValueInstance(key)
		if err != nil {
			return reflect.Value{}, err
		}

		valueCopy, err := duplicateValueInstance(value)
		if err != nil {
			return reflect.Value{}, err
		}

		newMap.SetMapIndex(keyCopy, valueCopy)
	}

	return newMap, nil
}

func duplicateValueForStruct(data reflect.Value) (reflect.Value, error) {
	newStruct := NewValueOfValue(data)

	for i := 0; i < data.NumField(); i++ {
		fieldValue := data.Field(i)
		fieldCopy, err := duplicateValueInstance(fieldValue)
		if err != nil {
			return reflect.Value{}, err
		}

		field := newStruct.Field(i)
		UnsafeValueSet(field, fieldCopy)
	}

	return newStruct, nil
}

func duplicateForBool(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	result.SetBool(data.Bool())
	return result, nil
}

func duplicateForInt(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	result.SetInt(data.Int())
	return result, nil
}

func duplicateForUint(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	result.SetUint(data.Uint())
	return result, nil
}

func duplicateForFloat(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	result.SetFloat(data.Float())
	return result, nil
}

func duplicateForString(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	result.SetString(data.String())
	return result, nil
}

func duplicateForPointer(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfValue(data)
	UnsafeValueSet(result, data)
	return result, nil
}

func duplicateForInterface(data reflect.Value) (reflect.Value, error) {
	value := NewValueOfValue(data)
	src := data

	for src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	newValue, err := duplicateValueInstance(src)
	if err != nil {
		return reflect.Value{}, err
	}

	value.Set(newValue)
	return value, nil
}

func duplicateForNilValue(data reflect.Value) (reflect.Value, error) {
	result := NewValueOfType(data.Type())
	return result, nil
}

func duplicateValueInstance(data reflect.Value) (reflect.Value, error) {
	isUntypedNil, isTypedNil := IsNilValue(data)
	if isUntypedNil {
		return NewUntypedNil(), ErrUntypedNil

	} else if isTypedNil {
		return duplicateForNilValue(data)
	}

	var err error
	var value reflect.Value

	switch data.Kind() {
	case reflect.Array, reflect.Slice:
		value, err = duplicateValueForArray(data)

	case reflect.Map:
		value, err = duplicateValueForMap(data)

	case reflect.Struct:
		value, err = duplicateValueForStruct(data)

	case reflect.Interface:
		value, err = duplicateForInterface(data)

	case reflect.Bool:
		value, err = duplicateForBool(data)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err = duplicateForInt(data)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value, err = duplicateForUint(data)

	case reflect.Float32, reflect.Float64:
		value, err = duplicateForFloat(data)

	case reflect.Ptr:
		value, err = duplicateForPointer(data)

	case reflect.String:
		value, err = duplicateForString(data)

	default:
		err = NewNotDuplicatableError(data.Kind())
	}

	return value, err
}

func DuplicateValueInstance(value reflect.Value) (reflect.Value, error) {
	valueInstance := ValueInstanceOf(value)
	return duplicateValueInstance(valueInstance)
}

func Duplicate(data interface{}) (interface{}, error) {
	dataValue := reflect.ValueOf(data)
	if !dataValue.IsValid() {
		return nil, ErrUntypedNil
	}

	instanceValue, chain := ValueToInstance(dataValue)
	copy, err := duplicateValueInstance(instanceValue)
	if err != nil {
		return nil, err
	}

	source := OriginOf(copy, chain)
	return source.Interface(), nil
}
