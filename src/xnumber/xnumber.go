package xnumber

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/radchukd/go-xschema/src/helpers"
)

type XNumber struct {
	validations map[string]helpers.XValidation[int]
}

func Create() XNumber {
	xn := XNumber{}
	xn.validations = make(map[string]helpers.XValidation[int])
	return xn
}

func FromTags(validationTags []string) XNumber {
	xn := Create()

	for _, v := range validationTags {
		nameArg := strings.Split(v, "=")

		switch nameArg[0] {
		case "Required":
			xn = xn.Required()
		case "Gt":
			n, _ := strconv.Atoi(nameArg[1])
			xn = xn.Gt(n)
		case "Gte":
			n, _ := strconv.Atoi(nameArg[1])
			xn = xn.Gte(n)
		case "Lt":
			n, _ := strconv.Atoi(nameArg[1])
			xn = xn.Lt(n)
		case "Lte":
			n, _ := strconv.Atoi(nameArg[1])
			xn = xn.Lte(n)
		case "MultipleOf":
			n, _ := strconv.Atoi(nameArg[1])
			xn = xn.MultipleOf(n)
		case "OneOf":
			var ooArr []int
			json.Unmarshal([]byte(nameArg[1]), &ooArr)
			xn = xn.OneOf(ooArr)
		}
	}

	return xn
}

func (xn XNumber) addValidation(ruleName string, err error, validation func(int) bool) XNumber {
	xn.validations[ruleName] = helpers.XValidation[int]{E: err, F: validation}
	return xn
}

func (xn XNumber) Validate(val interface{}) (bool, []error) {
	validationErrors := make([]error, 0)

	if reflect.TypeOf(val).Kind() == reflect.Float32 {
		return xn.Validate(int(val.(float32)))
	} else if reflect.TypeOf(val).Kind() == reflect.Float64 {
		return xn.Validate(int(val.(float64)))
	} else if reflect.TypeOf(val).Kind() != reflect.Int {
		validationErrors = append(validationErrors, fmt.Errorf("invalid type"))
		return false, validationErrors
	}

	value := val.(int)

	for _, validation := range xn.validations {
		isValid := validation.F(value)

		if !isValid {
			validationErrors = append(validationErrors, validation.E)
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func (xn XNumber) String() string {
	out := "XNumber("

	for validationName := range xn.validations {
		out += validationName + ","
	}

	out += ")"

	return out
}

func (xn XNumber) Required(errorMessage ...string) XNumber {
	return xn.addValidation(
		"Required",
		errors.New(append(errorMessage, "must be non-zero")[0]),
		func(value int) bool {
			return value != 0
		})
}

func (xn XNumber) Gt(gtValue int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"Gt",
		errors.New(append(errorMessage, fmt.Sprintf("must be greater than: %v", gtValue))[0]),
		func(value int) bool {
			return value > gtValue
		})
}

func (xn XNumber) Gte(gteValue int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"Gte",
		errors.New(append(errorMessage, fmt.Sprintf("must be greater or equal to: %v", gteValue))[0]),
		func(value int) bool {
			return value >= gteValue
		})
}

func (xn XNumber) Lt(ltValue int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"Lt",
		errors.New(append(errorMessage, fmt.Sprintf("must be lesser than: %v", ltValue))[0]),
		func(value int) bool {
			return value < ltValue
		})
}

func (xn XNumber) Lte(lteValue int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"Lte",
		errors.New(append(errorMessage, fmt.Sprintf("must be lesser or equal to: %v", lteValue))[0]),
		func(value int) bool {
			return value <= lteValue
		})
}

func (xn XNumber) MultipleOf(mtValue int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"MultipleOf",
		errors.New(append(errorMessage, fmt.Sprintf("must be a multiple of: %v", mtValue))[0]),
		func(value int) bool {
			return value%mtValue == 0
		})
}

func (xn XNumber) OneOf(possibleValues []int, errorMessage ...string) XNumber {
	return xn.addValidation(
		"OneOf",
		errors.New(append(errorMessage, fmt.Sprintf("must be one of: %v", possibleValues))[0]),
		func(value int) bool {
			for _, v := range possibleValues {
				if value == v {
					return true
				}
			}
			return false
		})
}
