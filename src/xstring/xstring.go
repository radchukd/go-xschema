package xstring

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/radchukd/go-xschema/src/helpers"
)

type XString struct {
	validations map[string]helpers.XValidation[string]
}

func Create() XString {
	xs := XString{}
	xs.validations = make(map[string]helpers.XValidation[string])
	return xs
}

func FromTags(validationTags []string) XString {
	xs := Create()

	for _, v := range validationTags {
		nameArg := strings.Split(v, "=")

		switch nameArg[0] {
		case "Required":
			xs = xs.Required()
		case "Alphanum":
			xs = xs.Alphanum()
		case "StartsWith":
			xs = xs.StartsWith(nameArg[1])
		case "EndsWith":
			xs = xs.EndsWith(nameArg[1])
		case "Lower":
			xs = xs.Lower()
		case "Upper":
			xs = xs.Upper()
		case "Length":
			ln, _ := strconv.Atoi(nameArg[1])
			xs = xs.Length(ln)
		case "Min":
			ln, _ := strconv.Atoi(nameArg[1])
			xs = xs.Min(ln)
		case "Max":
			ln, _ := strconv.Atoi(nameArg[1])
			xs = xs.Max(ln)
		case "Pattern":
			pt := regexp.MustCompile(nameArg[1])
			xs = xs.Pattern(*pt)
		case "Email":
			xs = xs.Email()
		case "URL":
			xs = xs.URL()
		case "UUID":
			xs = xs.UUID()
		case "OneOf":
			var ooArr []string
			json.Unmarshal([]byte(nameArg[1]), &ooArr)
			xs = xs.OneOf(ooArr)
		}
	}

	return xs
}

func (xs XString) addValidation(ruleName string, err error, validation func(string) bool) XString {
	xs.validations[ruleName] = helpers.XValidation[string]{E: err, F: validation}
	return xs
}

func (xs XString) Validate(val interface{}) (bool, []error) {
	validationErrors := make([]error, 0)

	if reflect.TypeOf(val).Kind() != reflect.String {
		validationErrors = append(validationErrors, fmt.Errorf("invalid type"))
		return false, validationErrors
	}

	value := val.(string)

	for _, validation := range xs.validations {
		isValid := validation.F(value)

		if !isValid {
			validationErrors = append(validationErrors, validation.E)
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func (xs XString) String() string {
	out := "XString("

	for validationName := range xs.validations {
		out += validationName + ","
	}

	out += ")"

	return out
}

func (xs XString) Required(errorMessage ...string) XString {
	return xs.addValidation(
		"Required",
		errors.New(append(errorMessage, "must be non-empty")[0]),
		func(value string) bool {
			return value != ""
		})
}

func (xs XString) Alphanum(errorMessage ...string) XString {
	return xs.addValidation(
		"Alphanum",
		errors.New(append(errorMessage, "must be alphanumeric")[0]),
		func(value string) bool {
			return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value)
		})
}

func (xs XString) StartsWith(prefix string, errorMessage ...string) XString {
	return xs.addValidation(
		"StartsWith",
		errors.New(append(errorMessage, fmt.Sprintf("must start with: %v", prefix))[0]),
		func(value string) bool {
			return strings.HasPrefix(value, prefix)
		})
}

func (xs XString) EndsWith(suffix string, errorMessage ...string) XString {
	return xs.addValidation(
		"EndsWith",
		errors.New(append(errorMessage, fmt.Sprintf("must start with: %v", suffix))[0]),
		func(value string) bool {
			return strings.HasSuffix(value, suffix)
		})
}

func (xs XString) Lower(errorMessage ...string) XString {
	return xs.addValidation(
		"Lower",
		errors.New(append(errorMessage, "must be lowercase")[0]),
		func(value string) bool {
			return regexp.MustCompile(`^[a-z]+$`).MatchString(value)
		})
}

func (xs XString) Upper(errorMessage ...string) XString {
	return xs.addValidation(
		"Upper",
		errors.New(append(errorMessage, "must be uppercase")[0]),
		func(value string) bool {
			return regexp.MustCompile(`^[A-Z]+$`).MatchString(value)
		})
}

func (xs XString) Length(length int, errorMessage ...string) XString {
	return xs.addValidation(
		"Length",
		errors.New(append(errorMessage, fmt.Sprintf("must be of length equal to: %v", length))[0]),
		func(value string) bool {
			return len(value) == length
		})
}

func (xs XString) Min(minLength int, errorMessage ...string) XString {
	return xs.addValidation(
		"Min",
		errors.New(append(errorMessage, fmt.Sprintf("must be of length greater than: %v", minLength))[0]),
		func(value string) bool {
			return len(value) >= minLength
		})
}

func (xs XString) Max(maxLength int, errorMessage ...string) XString {
	return xs.addValidation(
		"Max",
		errors.New(append(errorMessage, fmt.Sprintf("must be of length smaller than: %v", maxLength))[0]),
		func(value string) bool {
			return len(value) <= maxLength
		})
}

func (xs XString) Pattern(pattern regexp.Regexp, errorMessage ...string) XString {
	return xs.addValidation(
		"Pattern",
		errors.New(append(errorMessage, fmt.Sprintf("must match pattern: %s", pattern.String()))[0]),
		func(value string) bool {
			return pattern.MatchString(value)
		})
}

func (xs XString) Email(errorMessage ...string) XString {
	pattern := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return xs.Pattern(*pattern, errorMessage...)
}

func (xs XString) URL(errorMessage ...string) XString {
	pattern := regexp.MustCompile(`^(?:[a-zA-Z0-9]{1,62}(?:[-\.][a-zA-Z0-9]{1,62})+)(:\d+)?$`)
	return xs.Pattern(*pattern, errorMessage...)
}

func (xs XString) UUID(errorMessage ...string) XString {
	pattern := regexp.MustCompile(`^(?:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}|00000000-0000-0000-0000-000000000000)$`)
	return xs.Pattern(*pattern, errorMessage...)
}

func (xs XString) OneOf(possibleValues []string, errorMessage ...string) XString {
	return xs.addValidation(
		"OneOf",
		errors.New(append(errorMessage, fmt.Sprintf("must be one of: %v", possibleValues))[0]),
		func(value string) bool {
			for _, v := range possibleValues {
				if value == v {
					return true
				}
			}
			return false
		})
}
