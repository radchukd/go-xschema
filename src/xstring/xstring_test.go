package xstring_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/radchukd/go-xschema/src/xstring"
)

func TestValidate(t *testing.T) {
	xs := xstring.Create()

	value := ""

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Validate(%s) -> false; want true", value)
	}

	intValue := 123

	if isValid, _ := xs.Validate(intValue); isValid {
		t.Errorf("Validate(%v) -> true; want false", intValue)
	}
}

func TestRequired(t *testing.T) {
	var value string
	xs := xstring.Create().Required()

	value = "value"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Required(%s) -> false; want true", value)
	}

	value = ""

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Required(%s) -> true; want false", value)
	}
}

func TestCustomMessage(t *testing.T) {
	value := ""
	msg := "req"
	xs := xstring.Create().Required(msg)

	if _, err := xs.Validate(value); fmt.Sprint(err[0]) != msg {
		t.Errorf("Required(%s, %s) -> \"must be non-empty\"; want %s", value, msg, msg)
	}
}

func TestAlphanum(t *testing.T) {
	var value string
	xs := xstring.Create().Required().Alphanum()

	value = "abc"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Alphanum(%s) -> false; want true", value)
	}

	value = "123"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Alphanum(%s) -> false; want true", value)
	}

	value = "?_a1"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Alphanum(%s) -> true; want false", value)
	}
}

func TestStartsWith(t *testing.T) {
	var value string
	swValue := "A"
	xs := xstring.Create().StartsWith(swValue)

	value = "abc"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("StartsWith(%s,%s) -> true; want false", swValue, value)
	}

	value = "Abc"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("StartsWith(%s,%s) -> false; want true", swValue, value)
	}
}

func TestEndsWith(t *testing.T) {
	var value string
	ewValue := "C"
	xs := xstring.Create().EndsWith(ewValue)

	value = "abc"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("EndsWith(%s,%s) -> true; want false", ewValue, value)
	}

	value = "abC"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("EndsWith(%s,%s) -> false; want true", ewValue, value)
	}
}

func TestLower(t *testing.T) {
	var value string
	xs := xstring.Create().Lower()

	value = "Abc"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Lower(%s) -> true; want false", value)
	}

	value = "abc"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Lower(%s) -> false; want true", value)
	}
}

func TestUpper(t *testing.T) {
	var value string
	xs := xstring.Create().Upper()

	value = "Abc"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Upper(%s) -> true; want false", value)
	}

	value = "ABC"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Upper(%s) -> false; want true", value)
	}
}

func TestLength(t *testing.T) {
	var value string
	ln := 4
	xs := xstring.Create().Length(ln)

	value = "abc"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Length(%v,%s) -> true; want false", ln, value)
	}

	value = "abcd"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Length(%v,%s) -> false; want true", ln, value)
	}
}

func TestMin(t *testing.T) {
	var value string
	ln := 3
	xs := xstring.Create().Min(ln)

	value = "ab"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Min(%v,%s) -> true; want false", ln, value)
	}

	value = "abcd"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Min(%v,%s) -> false; want true", ln, value)
	}
}

func TestMax(t *testing.T) {
	var value string
	ln := 3
	xs := xstring.Create().Max(ln)

	value = "abcd"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Max(%v,%s) -> true; want false", ln, value)
	}

	value = "ab"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Max(%v,%s) -> false; want true", ln, value)
	}
}

func TestPattern(t *testing.T) {
	var value string
	pattern := regexp.MustCompile(`([A-Z])\w+`)
	xs := xstring.Create().Pattern(*pattern)

	value = "abcd"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Pattern(%s,%s) -> true; want false", pattern.String(), value)
	}

	value = "RegExr"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Pattern(%s,%s) -> false; want true", pattern.String(), value)
	}
}

func TestEmail(t *testing.T) {
	var value string
	xs := xstring.Create().Email()

	value = "email"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("Email(%s) -> true; want false", value)
	}

	value = "email@example.com"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("Email(%s) -> false; want true", value)
	}
}

func TestURL(t *testing.T) {
	var value string
	xs := xstring.Create().URL()

	value = "url"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("URL(%s) -> true; want false", value)
	}

	value = "example.com"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("URL(%s) -> false; want true", value)
	}
}

func TestUUID(t *testing.T) {
	var value string
	xs := xstring.Create().UUID()

	value = "uuid"

	if isValid, _ := xs.Validate(value); isValid {
		t.Errorf("UUID(%s) -> true; want false", value)
	}

	value = "2f4e6b64-557a-4f07-ab04-ec31a7d9f5e0"

	if isValid, _ := xs.Validate(value); !isValid {
		t.Errorf("UUID(%s) -> false; want true", value)
	}
}

func TestOneOf(t *testing.T) {
	var value string
	possibleValues := []string{"apple", "mango", "banana"}
	xn := xstring.Create().OneOf(possibleValues)

	value = "mango"

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("OneOf(%v,%v) -> false; want true", possibleValues, value)
	}

	value = "orange"

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("OneOf(%v,%v) -> true; want false", possibleValues, value)
	}

	value = "apple"
	xn = xstring.Create().OneOf(possibleValues[1:])

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("OneOf(%v,%v) -> true; want false", possibleValues, value)
	}
}
