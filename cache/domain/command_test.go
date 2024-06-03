package domain

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input    []byte
		expected *Command
		err      error
	}{
		// SET
		{
			input:    []byte("SET key value"),
			expected: &Command{Type: CMDSet, Key: "key", Value: "value"},
			err:      nil,
		},
		{
			input:    []byte("SET key"), // wrong arguments for SET command
			expected: nil,
			err:      errors.New("wrong arguments for SET command"),
		},
		{
			input:    []byte("INVALID key value"), // Invalid command for SET
			expected: nil,
			err:      errors.New("invalid command"),
		},
		{
			input:    []byte("SET key value extra"), // wrong arguments for SET command
			expected: nil,
			err:      errors.New("wrong arguments for SET command"),
		},
		// GET
		{
			input:    []byte("GET key"),
			expected: &Command{Type: CMDGet, Key: "key"},
			err:      nil,
		},
		{
			input:    []byte("GET"), // wrong arguments for GET command
			expected: nil,
			err:      errors.New("wrong arguments for GET command"),
		},
		{
			input:    []byte("INVALID key"), // invalid command for GET
			expected: nil,
			err:      errors.New("invalid command"),
		},
	}

	for _, tt := range tests {
		result, err := ParseCommand(tt.input)
		if !reflect.DeepEqual(result, tt.expected) || !equalError(err, tt.err) {
			t.Errorf("ParseCommand(%s) = (%v, %v), expected (%v, %v)", tt.input, result, err, tt.expected, tt.err)
		}
	}
}

func equalError(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}
	return err1.Error() == err2.Error()
}
