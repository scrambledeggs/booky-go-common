package helpers

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type Schema struct {
	Body            *map[string]SchemaType
	QueryParameters *map[string]SchemaType
	PathParameters  *map[string]SchemaType
}

type SchemaTypes string

const (
	Number  SchemaTypes = "number"
	String  SchemaTypes = "string"
	Boolean SchemaTypes = "boolean"
	UUID    SchemaTypes = "uuid"
)

type SchemaType struct {
	Type     SchemaTypes
	Required bool
}

// Validates request payload base on the given schema
func ValidatePayload(request events.APIGatewayProxyRequest, schema Schema) error {
	t := reflect.TypeOf(schema)

	errCh := make(chan error, t.NumField())
	var wg sync.WaitGroup

	for i := range t.NumField() {
		value := t.Field(i)
		wg.Add(1)
		go func(value reflect.StructField) {
			defer wg.Done()
			switch value.Name {
			case "Body":
				err := validateBody(request.Body, schema.Body)
				if err != nil {
					errCh <- err
				}
			case "QueryParameters":
				err := validateQueryStringParameters(request.QueryStringParameters, schema.QueryParameters)
				if err != nil {
					errCh <- err
				}
			case "PathParameters":
				err := validatePathParameters(request.PathParameters, schema.PathParameters)
				if err != nil {
					errCh <- err
				}
			}
		}(value)
	}

	// Wait for all goroutines to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Return the first error encountered, if any
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func validateQueryStringParameters(queryStringParameters map[string]string, schema *map[string]SchemaType) error {
	if schema == nil {
		return nil
	}

	for key, value := range *schema {
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)

		for i := range t.NumField() {
			field := t.Field(i).Name
			value := v.Field(i).Interface()

			err := validate(queryStringParameters[key], value, field)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validatePathParameters(pathParameters map[string]string, schema *map[string]SchemaType) error {
	if schema == nil {
		return nil
	}

	for key, value := range *schema {
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)

		for i := range t.NumField() {
			field := t.Field(i).Name
			value := v.Field(i).Interface()

			err := validate(pathParameters[key], value, field)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateBody(body string, schema *map[string]SchemaType) error {
	if schema == nil {
		return nil
	}

	var bodyObj map[string]any

	err := json.Unmarshal([]byte(body), &bodyObj)

	if err != nil {
		return errors.New("body is not a valid JSON")
	}

	for key, value := range *schema {
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)

		for i := range t.NumField() {
			field := t.Field(i).Name
			value := v.Field(i).Interface()

			err := validate(bodyObj[key], value, field)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validate(payload any, value any, field string) error {
	switch field {
	case "Required":
		if value.(bool) && (payload == nil || payload == "") {
			return errors.New("payload is required")
		}
		return nil
	case "Type":
		if payload == nil || payload == "" {
			return nil
		}
		switch value.(SchemaTypes) {
		case Number:
			return validateNumber(payload.(string))
		case String:
			return nil
		case Boolean:
			return validateBoolean(payload.(string))
		case UUID:
			return validateUUID(payload.(string))
		default:
			return errors.New("invalid schema type")
		}
	default:
		return nil
	}
}

func validateNumber(payload string) error {
	_, err := strconv.ParseFloat(payload, 64)

	if err != nil {
		return errors.New("payload is not a number")
	}

	return nil
}

func validateBoolean(payload string) error {
	_, err := strconv.ParseBool(payload)

	if err != nil {
		return errors.New("payload is not a boolean")
	}

	return nil
}

func validateUUID(payload string) error {
	_, err := uuid.Parse(payload)

	if err != nil {
		return errors.New("payload is not a uuid")
	}

	return nil
}
