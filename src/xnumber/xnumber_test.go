package xnumber_test

import (
	"fmt"
	"testing"

	"github.com/radchukd/go-xschema/src/xnumber"
)

func TestValidate(t *testing.T) {
	xn := xnumber.Create()

	intValue := 123

	if isValid, _ := xn.Validate(intValue); !isValid {
		t.Errorf("Validate(%v) -> false; want true", intValue)
	}

	floatValue := 123.0

	if isValid, _ := xn.Validate(floatValue); !isValid {
		t.Errorf("Validate(%v) -> false; want true", floatValue)
	}

	value := ""

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Validate(%s) -> true; want false", value)
	}
}

func TestRequired(t *testing.T) {
	var value int
	xn := xnumber.Create().Required()

	value = 123

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("Required(%v) -> false; want true", value)
	}

	value = 0

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Required(%v) -> true; want false", value)
	}
}

func TestCustomMessage(t *testing.T) {
	value := 0
	msg := "req"
	xn := xnumber.Create().Required(msg)

	if _, err := xn.Validate(value); fmt.Sprint(err[0]) != msg {
		t.Errorf("Required(%v, %s) -> \"must be non-empty\"; want %s", value, msg, msg)
	}
}

func TestGt(t *testing.T) {
	var value int
	compVal := 0
	xn := xnumber.Create().Gt(compVal)

	value = -123

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Gt(%v,%v) -> true; want false", compVal, value)
	}

	value = 123

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("Gt(%v,%v) -> false; want true", compVal, value)
	}
}

func TestGte(t *testing.T) {
	var value int
	compVal := 0
	xn := xnumber.Create().Gte(compVal)

	value = -123

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Gte(%v,%v) -> true; want false", compVal, value)
	}

	value = 0

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("Gte(%v,%v) -> false; want true", compVal, value)
	}
}

func TestLt(t *testing.T) {
	var value int
	compVal := 0
	xn := xnumber.Create().Lt(compVal)

	value = 123

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Lt(%v,%v) -> true; want false", compVal, value)
	}

	value = -123

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("Lt(%v,%v) -> false; want true", compVal, value)
	}
}

func TestLte(t *testing.T) {
	var value int
	compVal := 0
	xn := xnumber.Create().Lte(compVal)

	value = 123

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("Lte(%v,%v) -> true; want false", compVal, value)
	}

	value = 0

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("Lte(%v,%v) -> false; want true", compVal, value)
	}
}

func TestMultipleOf(t *testing.T) {
	var value int
	mValue := 3
	xn := xnumber.Create().MultipleOf(mValue)

	value = 17

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("MultipleOf(%v,%v) -> true; want false", mValue, value)
	}

	value = 27

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("MultipleOf(%v,%v) -> false; want true", mValue, value)
	}
}

func TestOneOf(t *testing.T) {
	var value int
	possibleValues := []int{7, 17, 27}
	xn := xnumber.Create().OneOf(possibleValues)

	value = 17

	if isValid, _ := xn.Validate(value); !isValid {
		t.Errorf("OneOf(%v,%v) -> false; want true", possibleValues, value)
	}

	value = 9

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("OneOf(%v,%v) -> true; want false", possibleValues, value)
	}

	value = 7
	xn = xnumber.Create().OneOf(possibleValues[1:])

	if isValid, _ := xn.Validate(value); isValid {
		t.Errorf("OneOf(%v,%v) -> true; want false", possibleValues, value)
	}
}
