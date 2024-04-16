package validator_test

import (
	"testing"

	"github.com/andremfp/snippetbox/internal/validator"
)

func TestAddFieldError(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		message  string
		expected string
	}{
		{
			name:     "Adding new key that does not exists",
			key:      "key1",
			message:  "message1",
			expected: "message1",
		},
		{
			name:     "Adding key that already exists",
			key:      "key1",
			message:  "message2",
			expected: "message1",
		},
		{
			name:     "Adding key to empty map",
			key:      "key2",
			message:  "message3",
			expected: "message3",
		},
	}

	v := &validator.Validator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.AddFieldError(tt.key, tt.message)
			if v.FieldErrors[tt.key] != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, v.FieldErrors[tt.key])
			}
		})
	}
}

func TestCheckField(t *testing.T) {
	tests := []struct {
		name      string
		ok        bool
		key       string
		message   string
		numErrors int
	}{
		{
			name:      "Need to add an error",
			ok:        true,
			numErrors: 0,
		},
		{
			name:      "No need to add an error",
			ok:        false,
			numErrors: 1,
		},
	}

	v := &validator.Validator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.CheckField(tt.ok, "key1", "message1")
			if len(v.FieldErrors) != tt.numErrors {
				t.Errorf("Expected %v errors, but got %v errors", tt.numErrors, len(v.FieldErrors))
			}
		})
	}
}

func TestPermittedValue(t *testing.T) {

	tests := []struct {
		name            string
		permittedValues []string
		value           string
		expected        bool
	}{
		{
			name:            "Value in permitted values list",
			permittedValues: []string{"1", "7", "365"},
			value:           "1",
			expected:        true,
		},
		{
			name:            "Value not in permitted values list",
			permittedValues: []string{"1", "7", "365"},
			value:           "3",
			expected:        false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := validator.PermittedValue(tt.value, tt.permittedValues...); got != tt.expected {
				t.Errorf("PermittedValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}
