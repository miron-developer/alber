package orm

import (
	"errors"
	"reflect"
)

// MakeArrFromStruct struct -> []interface
func MakeArrFromStruct(data interface{}) []interface{} {
	arr := []interface{}{}
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		arr = append(arr, v.Field(i).Interface())
	}
	return arr
}

// FillStructFromMap fill the current struct from map
func FillStructFromMap(sampleStruct interface{}, data map[string]interface{}) error {
	structValue := reflect.ValueOf(sampleStruct).Elem()

	for k, v := range data {
		structFieldValue := structValue.FieldByName(k)

		if !structFieldValue.IsValid() {
			return errors.New("no such field: " + k + " in obj")
		}
		if !structFieldValue.CanSet() {
			return errors.New("can't set " + k + " field value")
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(v)
		if structFieldType != val.Type() {
			return errors.New("provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}
	return nil
}

// MapFromStructAndMatrix [][]interface{}{} + Struct sample -> []map[string]interface{}{}
func MapFromStructAndMatrix(data [][]interface{}, sampleStruct interface{}, additionalFields ...string) []map[string]interface{} {
	structLen := reflect.ValueOf(sampleStruct).NumField()
	t := reflect.TypeOf(sampleStruct)
	result := []map[string]interface{}{}

	for _, currentRow := range data {
		oneRow := map[string]interface{}{}
		for i := 0; i < structLen; i++ {
			jsonName := t.Field(i).Tag.Get("json")
			if jsonName == "" {
				continue
			}
			oneRow[jsonName] = currentRow[i]
		}

		// add additional fields for result
		for i, v := range additionalFields {
			oneRow[v] = currentRow[i+structLen]
		}
		result = append(result, oneRow)
	}
	return result
}

// FromINT64ToINT convert int64 -> int
func FromINT64ToINT(number interface{}) int {
	return int(number.(int64))
}
