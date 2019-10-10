package marshalling

import (
	"reflect"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// CustomMarshalMap is a custom marshal func which marshals Go value to AttributeValues
func CustomMarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av := dynamodbattribute.NewEncoder()
	av.NullEmptyString = false // keep empty strings

	// convert interface{} to map[string]*dynamodb.AttributeValue
	aav, err := av.Encode(in)
	if err != nil || av == nil || aav.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	checkAndConvertIfNull(aav.M)
	return aav.M, nil
}

func checkAndConvertIfNull(attr map[string]*dynamodb.AttributeValue) {
	// loop over map[string]*dynamodb.AttributeValue
	for key, elem := range attr {
		// convert elem to Go value from *dynamodb.AttributeValue
		var i interface{}
		_ = dynamodbattribute.Unmarshal(elem, &i)

		// change go value to reflect.Value for null checking
		r := reflect.ValueOf(i)
		if emptyValue(r) {
			delete(attr, key)

			// convert back to *dynamodb.AttributeValue
			t := true
			attr[key] = &dynamodb.AttributeValue{NULL: &t}
		}

		if r.Kind() == reflect.Map {
			// recursively check maps for empty values
			checkAndConvertIfNull(elem.M)
		}
	}
}

func emptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
