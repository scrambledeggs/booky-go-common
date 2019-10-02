package custommarshal

import (
	"reflect"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Map is a custom marshal func which marshals Go value to AttributeValues
func Map(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av := dynamodbattribute.NewEncoder()
	av.NullEmptyString = false // keep empty strings

	// convert interface{} to map[string]*dynamodb.AttributeValue
	aav, err := av.Encode(in)
	if err != nil || av == nil || aav.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	// loop over map[string]*dynamodb.AttributeValue
	for key, elem := range aav.M {
		// convert elem to Go value from *dynamodb.AttributeValue
		var i interface{}
		_ = dynamodbattribute.Unmarshal(elem, &i)

		// change go value to reflect.Value for null checking
		r := reflect.ValueOf(i)
		if emptyValue(r) {
			delete(aav.M, key)

			// convert back to *dynamodb.AttributeValue
			t := true
			aav.M[key] = &dynamodb.AttributeValue{NULL: &t}
		}
	}

	return aav.M, nil
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
