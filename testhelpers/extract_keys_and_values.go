package testhelpers

import "reflect"

// Extract keys and values from mock struct to add record to mock DB
func ExtractKeysAndValues(obj any) ([]string, []interface{}) {
	var keys []string
	var values []interface{}

	objKeys := reflect.TypeOf(obj)
	objValues := reflect.ValueOf(obj)

	for i := 0; i < objValues.NumField(); i++ {
		keys = append(keys, objKeys.Field(i).Name)
		values = append(values, objValues.Field(i).Interface())
	}

	return keys, values
}
