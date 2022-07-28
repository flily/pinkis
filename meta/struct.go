package meta

import (
	"reflect"
)

func GetFieldValue(data interface{}, field string) (reflect.Value, error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	var err error
	fieldValue := dataValue.FieldByName(field)
	if !fieldValue.IsValid() {
		err = NewMetaError("no field '%s'", field)
	}

	return fieldValue, err
}

func GetField(data interface{}, field string) (interface{}, error) {
	fieldValue, err := GetFieldValue(data, field)
	if err != nil {
		return nil, err
	}

	return fieldValue.Interface(), nil
}

func SetField(data interface{}, field string, value interface{}) (interface{}, error) {
	fieldValue, err := GetFieldValue(data, field)
	if err != nil {
		return nil, err
	}

	valueValue := reflect.ValueOf(value)
	if !fieldValue.CanSet() {
		return nil, NewMetaError("field '%s' can not be set", field)
	}

	fieldValueType := fieldValue.Type()
	if valueValue.Type() != fieldValueType {
		if valueValue.CanConvert(fieldValueType) {
			convertedValue := valueValue.Convert(fieldValueType)
			fieldValue.Set(convertedValue)

		} else {
			return nil, NewMetaError("field '%s' requires type %s, but %s",
				field, fieldValue.Type(), valueValue.Type())
		}

	} else {
		fieldValue.Set(valueValue)
	}

	return fieldValue.Interface(), nil
}
