package xschema_test

import (
	"regexp"
	"testing"

	"github.com/radchukd/go-xschema/src/xnumber"
	"github.com/radchukd/go-xschema/src/xschema"
	"github.com/radchukd/go-xschema/src/xstring"
)

func TestClone(t *testing.T) {}

func TestMerge(t *testing.T) {}

func TestValidateKey(t *testing.T) {
	var value string
	key := "Name"
	schema := xschema.Create().
		AddString("Name", xstring.Create().Required().Alphanum().Min(2))

	value = "abc12"

	if isValid, _ := schema.ValidateKey(key, value); !isValid {
		t.Errorf("ValidateKey(%s, %v) -> false; want true", key, value)
	}

	value = "_______"

	if isValid, _ := schema.ValidateKey("Name", value); isValid {
		t.Errorf("ValidateKey(%s, %v) -> false; want true", key, value)
	}
}

func TestSValidateKey(t *testing.T) {
	var value string
	key := "Name"
	schema := xschema.Create().
		AddString("Name", xstring.Create().Required().Alphanum().Min(2))

	value = "John"

	if isValid, _ := schema.SValidateKey(key, value); !isValid {
		t.Errorf("SValidateKey(%s, %v) -> false; want true", key, value)
	}

	key = "Age"

	if isValid, _ := schema.SValidateKey(key, value); isValid {
		t.Errorf("SValidateKey(%s, %v) -> true; want false", key, value)
	}
}

func TestValidateMap(t *testing.T) {
	var values map[string]interface{}
	schema := xschema.Create().
		AddString("FirstName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`))).
		AddString("LastName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`))).
		AddNumber("Age", xnumber.Create().Gte(18))

	values = make(map[string]interface{})

	values["FirstName"] = "John"
	values["LastName"] = "doe"
	values["Age"] = 0

	if isValid, _ := schema.ValidateMap(values); isValid {
		t.Errorf("ValidateMap(%v) -> true; want false", values)
	}

	values["LastName"] = "Doe"
	values["Age"] = 18

	if isValid, _ := schema.ValidateMap(values); !isValid {
		t.Errorf("ValidateMap(%v) -> false; want true", values)
	}
}

func TestSValidateMap(t *testing.T) {
	var values map[string]interface{}
	schema := xschema.Create().
		AddString("FirstName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`)))

	values = make(map[string]interface{})

	values["FirstName"] = "John"
	values["LastName"] = "Doe"

	if isValid, _ := schema.SValidateMap(values); isValid {
		t.Errorf("SValidateMap(%v) -> true; want false", values)
	}

	delete(values, "LastName")

	if isValid, _ := schema.ValidateMap(values); !isValid {
		t.Errorf("SValidateMap(%v) -> false; want true", values)
	}
}

func TestValidateStruct(t *testing.T) {
	type User struct {
		FirstName string
		LastName  string
		Age       uint8
	}

	var value User
	schema := xschema.Create().
		AddString("FirstName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`))).
		AddNumber("Age", xnumber.Create().Gte(18))

	value = User{"John", "Doe", 0}

	if isValid, _ := schema.ValidateStruct(value); isValid {
		t.Errorf("ValidateStruct(%v) -> true; want false", value)
	}

	value.Age = 18

	if isValid, _ := schema.ValidateStruct(value); !isValid {
		t.Errorf("ValidateStruct(%v) -> false; want true", value)
	}
}

func TestSValidateStruct(t *testing.T) {
	type User struct {
		FirstName string
		LastName  string
		Age       uint8
	}

	var value User
	schema := xschema.Create().
		AddString("FirstName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`))).
		AddNumber("Age", xnumber.Create().Gte(18))

	value = User{"John", "Doe", 18}

	if isValid, _ := schema.SValidateStruct(value); isValid {
		t.Errorf("SValidateStruct(%v) -> true; want false", value)
	}

	schema = schema.AddString("LastName", xstring.Create().Required().Pattern(*regexp.MustCompile(`^[A-Z]{1}[a-z]+$`)))

	if isValid, _ := schema.SValidateStruct(value); !isValid {
		t.Errorf("SValidateStruct(%v) -> false; want true", value)
	}
}

func TestValidateTaggedStruct(t *testing.T) {
	type User struct {
		FirstName string `x:"Required,Min=3,Pattern=^[A-Z]{1}[a-z]+$"`
		LastName  string `x:"Required,Min=3,Pattern=^[A-Z]{1}[a-z]+$"`
		Age       uint8  `x:"Required,Gte=18"`
	}

	value := User{"John", "doe", 0}

	if isValid, _ := xschema.ValidateTaggedStruct(value); isValid {
		t.Errorf("ValidateTaggedStruct(%v) -> true; want false", value)
	}

	value.LastName = "Doe"
	value.Age = 18

	if isValid, _ := xschema.ValidateTaggedStruct(value); !isValid {
		t.Errorf("ValidateTaggedStruct(%v) -> false; want true", value)
	}
}
