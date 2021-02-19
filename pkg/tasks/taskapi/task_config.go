package taskapi

import (
	"fmt"
	"reflect"
)

type propertyNotFoundError struct {
	name string
}

func (e propertyNotFoundError) Error() string {
	return fmt.Sprintf("property \"%s\" not found", e.name)
}

// TaskConfig represents a task as defined in dev.yml
type TaskConfig struct {
	name    string
	payload interface{}
}

func NewTaskConfig(definition interface{}) (*TaskConfig, error) {
	val := reflect.ValueOf(definition)

	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) != 1 {
			return nil, fmt.Errorf("invalid map length")
		}
		name, ok := keys[0].Interface().(string)
		if !ok {
			return nil, fmt.Errorf("task name should be a string")
		}
		payload := val.MapIndex(keys[0]).Interface()
		return &TaskConfig{name: name, payload: payload}, nil
	}

	if val.Kind() == reflect.String {
		return &TaskConfig{name: definition.(string), payload: nil}, nil
	}

	return nil, fmt.Errorf("invalid task: \"%+v\"", definition)
}

// IsHash returns a boolean indicating whether the payload is a hash
func (c *TaskConfig) IsHash() bool {
	return reflect.ValueOf(c.payload).Kind() == reflect.Map
}

// GetListOfStrings expects the payload to be a list of string, returns it.
//
// YAML Example:
//   - TASKNAME:
//     - foo
//     - bar
func (c *TaskConfig) GetListOfStrings() ([]string, error) {
	return asListOfStrings(c.payload)
}

// GetListOfStringsPropertyDefault expects the payload to be a list of string, returns it.
// Or returns the default value specified if the payload is nil.
//
// YAML Example:
//   - TASKNAME:
//     - foo
//     - bar
func (c *TaskConfig) GetListOfStringsPropertyDefault(name string, defaultValue []string) ([]string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return nil, err
	}
	return asListOfStrings(value)
}

// GetStringPropertyAllowSingle expects one of those two situations:
// 1. the payload is a string, returns it directly
// 2. the payload is a hash of string, returns the value for the name specified
//
// YAML Example:
//   - TASKNAME: foo
// or
//   - TASKNAME:
//       PROPERTYNAME: foo
func (c *TaskConfig) GetStringPropertyAllowSingle(name string) (string, error) {
	if value, ok := c.payload.(string); ok {
		return value, nil
	}
	return c.GetStringProperty(name)
}

// GetStringProperty expects the payload to be a hash of string, returns the value for the name specified.
//
// YAML Example:
//   - TASKNAME:
//       PROPERTYNAME: foo
func (c *TaskConfig) GetStringProperty(name string) (string, error) {
	value, err := c.getProperty(name)
	if err != nil {
		return "", err
	}
	str, err := asString(value)
	if err != nil {
		return "", fmt.Errorf(`key "%s": %w`, name, err)
	}
	return str, nil
}

// GetStringPropertyDefault expects the payload to be a hash of string, returns the value for the name specified.
// Or returns the default value specified if the payload is nil.
//
// YAML Example:
//   - TASKNAME:
//       PROPERTYNAME: foo
func (c *TaskConfig) GetStringPropertyDefault(name string, defaultValue string) (string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return "", err
	}
	str, err := asString(value)
	if err != nil {
		return "", fmt.Errorf(`key "%s": %w`, name, err)
	}
	return str, nil
}

// GetBooleanPropertyDefault expects the payload to be a hash, returns the value as a boolean for the name specified.
// Or returns the default value specified if the payload is nil.
//
// YAML Example:
//   - TASKNAME:
//       PROPERTYNAME: on
func (c *TaskConfig) GetBooleanPropertyDefault(name string, defaultValue bool) (bool, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return false, err
	}
	return asBool(value)
}

func (c *TaskConfig) getProperty(name string) (interface{}, error) {
	if c.payload == nil {
		return nil, propertyNotFoundError{name: name}
	}

	properties, ok := c.payload.(map[interface{}]interface{})
	if !ok {
		return "", fmt.Errorf("expecting a hash, found a %T (%+v)", c.payload, c.payload)
	}

	value, present := properties[name]
	if present {
		return value, nil
	}
	return nil, propertyNotFoundError{name: name}
}

func (c *TaskConfig) getPropertyDefault(name string, defaultValue interface{}) (interface{}, error) {
	value, err := c.getProperty(name)
	if err == nil {
		return value, nil
	}
	if _, ok := err.(propertyNotFoundError); ok {
		return defaultValue, nil
	}
	return nil, err
}
