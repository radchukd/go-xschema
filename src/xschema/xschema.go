package xschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/radchukd/go-xschema/src/helpers"
	"github.com/radchukd/go-xschema/src/xnumber"
	"github.com/radchukd/go-xschema/src/xstring"
)

var tagName = "x"

type XSchema struct {
	values map[string]helpers.XObject
}

func Create() XSchema {
	schema := XSchema{}
	schema.values = make(map[string]helpers.XObject)
	return schema
}

func Clone(schema XSchema) XSchema {
	newSchema := XSchema{}
	newSchema.values = make(map[string]helpers.XObject)

	for k, v := range schema.values {
		newSchema.values[k] = v
	}

	return newSchema
}

func Merge(s1 XSchema, s2 XSchema) XSchema {
	newSchema := XSchema{}
	newSchema.values = make(map[string]helpers.XObject)

	for k, v := range s1.values {
		newSchema.values[k] = v
	}

	for k, v := range s2.values {
		newSchema.values[k] = v
	}

	return newSchema
}

func (schema XSchema) AddString(key string, xs xstring.XString) XSchema {
	schema.values[key] = xs
	return schema
}

func (schema XSchema) AddNumber(key string, xn xnumber.XNumber) XSchema {
	schema.values[key] = xn
	return schema
}

func (schema XSchema) ValidateKey(schemaKey string, value interface{}) (bool, []error) {
	for key, val := range schema.values {
		if key == schemaKey {
			if isValid, errors := val.Validate(value); !isValid {
				return false, errors
			}

			return true, nil
		}
	}

	return true, nil
}

func (schema XSchema) SValidateKey(schemaKey string, value interface{}) (bool, []error) {
	for key, val := range schema.values {
		if key == schemaKey {
			if isValid, errors := val.Validate(value); !isValid {
				return false, errors
			}

			return true, nil
		}
	}

	errs := make([]error, 0)
	return false, append(errs, errors.New("invalid key"))
}

func (schema XSchema) ValidateMap(values map[string]interface{}) (bool, map[string][]error) {
	validationErrors := make(map[string][]error)

	for key, value := range values {
		if isValid, errors := schema.ValidateKey(key, value); !isValid {
			validationErrors[fmt.Sprintf("%s(%v)", key, value)] = errors
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func (schema XSchema) SValidateMap(values map[string]interface{}) (bool, map[string][]error) {
	validationErrors := make(map[string][]error)

	for key, value := range values {
		if isValid, errors := schema.SValidateKey(key, value); !isValid {
			validationErrors[fmt.Sprintf("%s(%v)", key, value)] = errors
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func (schema XSchema) ValidateStruct(obj interface{}) (bool, map[string][]error) {
	var mappedObj map[string]interface{}
	inrec, _ := json.Marshal(obj)
	json.Unmarshal(inrec, &mappedObj)

	return schema.ValidateMap(mappedObj)
}

func (schema XSchema) SValidateStruct(obj interface{}) (bool, map[string][]error) {
	var mappedObj map[string]interface{}
	inrec, _ := json.Marshal(obj)
	json.Unmarshal(inrec, &mappedObj)
	return schema.SValidateMap(mappedObj)
}

func ValidateTaggedStruct(obj interface{}) (bool, map[string][]error) {
	schema := Create()
	t := reflect.TypeOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		if tag == "" {
			continue
		}

		validationTags := strings.Split(tag, ",")

		switch field.Type.Kind() {
		case reflect.String:
			schema = schema.AddString(field.Name, xstring.FromTags(validationTags))
		default:
			schema = schema.AddNumber(field.Name, xnumber.FromTags(validationTags))
		}
	}

	return schema.ValidateStruct(obj)
}
