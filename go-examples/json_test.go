package go_examples

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	validator "github.com/go-playground/validator/v10"
)

var ErrParsing = errors.New("ErrParsing")

var validate *validator.Validate

type Fruit string

const (
	FruitApple  Fruit = "apple"
	FruitOrange Fruit = "orange"
	FruitBanana Fruit = "banana"
)

func FruitFromString(s string) (Fruit, error) {
	fruit := Fruit(s)

	switch fruit {
	case FruitApple:
		fallthrough
	case FruitOrange:
		fallthrough
	case FruitBanana:
		return fruit, nil
	}

	return fruit, fmt.Errorf("failed to parse Fruit from %s: %w", s, ErrParsing)
}

func (f *Fruit) Validate() error {
	_, err := FruitFromString(string(*f))

	return err
}

type JsonWithFruit struct {
	Fruit       Fruit   `json:"fruit" validate:"validateFn"`
	Owner       string  `json:"owner" validate:"required"`
	Description *string `json:"description,omitempty"`
}

func TestMain(m *testing.M) {
	validate = validator.New()
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestValidJson(t *testing.T) {
	originalJson := `{
  "fruit": "apple",
  "owner": "John",
  "description": "a sweet one"
}`

	var typedValue JsonWithFruit
	expectedDescription := "a sweet one"

	err := json.Unmarshal([]byte(originalJson), &typedValue)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %#v", err)
	}
	if typedValue.Fruit != FruitApple {
		t.Errorf("expected fruit [%s] but got [%s]", FruitApple, typedValue.Fruit)
	}
	if *typedValue.Description != expectedDescription {
		t.Errorf("expected description [%s] but got [%s]", expectedDescription, *typedValue.Description)
	}

	err = validate.Struct(&typedValue)
	if err != nil {
		t.Errorf("expected valid JSON, but got error %+v", err)
	}

	result, err := json.MarshalIndent(typedValue, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal JSON: %#v", err)
	}
	resultingJson := string(result)
	if originalJson != resultingJson {
		t.Errorf("expected: \n%s\nactual: \n%s\n", originalJson, resultingJson)
	}
}

func TestValidJsonWithoutOptionalField(t *testing.T) {
	originalJson := `{
  "fruit": "apple",
  "owner": "John"
}`

	var typedValue JsonWithFruit

	err := json.Unmarshal([]byte(originalJson), &typedValue)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %#v", err)
	}
	if typedValue.Fruit != FruitApple {
		t.Errorf("expected fruit [%s] but got [%s]", FruitApple, typedValue.Fruit)
	}
	if typedValue.Description != nil {
		t.Errorf("expected no description but got [%s]", *typedValue.Description)
	}

	err = validate.Struct(&typedValue)
	if err != nil {
		t.Errorf("expected valid JSON, but got error %+v", err)
	}

	result, err := json.MarshalIndent(typedValue, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal JSON: %#v", err)
	}
	resultingJson := string(result)
	if originalJson != resultingJson {
		t.Errorf("expected: \n%s\nactual: \n%s\n", originalJson, resultingJson)
	}
}

func TestInvalidJsonRequiredEnumFieldMissing(t *testing.T) {
	originalJson := `{
  "owner": "John"
}`

	var typedValue JsonWithFruit

	err := json.Unmarshal([]byte(originalJson), &typedValue)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %#v", err)
	}

	err = validate.Struct(&typedValue)

	if err == nil || !strings.Contains(err.Error(), "Error:Field validation for 'Fruit' failed") {
		t.Error("expected validation error but got none")
	}
}

func TestInvalidJsonRequiredStringFieldMissing(t *testing.T) {
	originalJson := `{
  "fruit": "apple"
}`

	var typedValue JsonWithFruit

	err := json.Unmarshal([]byte(originalJson), &typedValue)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %#v", err)
	}

	err = validate.Struct(&typedValue)

	if err == nil || !strings.Contains(err.Error(), "Error:Field validation for 'Owner' failed") {
		t.Error("expected validation error but got none")
	}
}

func TestInvalidJsonWrongEnumValue(t *testing.T) {
	originalJson := `{
  "fruit": "appleWithTypo",
  "owner": "John"
}`

	var typedValue JsonWithFruit

	err := json.Unmarshal([]byte(originalJson), &typedValue)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %#v", err)
	}

	err = validate.Struct(&typedValue)

	if err == nil || !strings.Contains(err.Error(), "Error:Field validation for 'Fruit' failed") {
		t.Error("expected validation error but got none")
	}
}
