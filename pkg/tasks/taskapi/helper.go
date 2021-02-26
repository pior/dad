package taskapi

import "fmt"

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	return "", fmt.Errorf("expecting a string, found a %T (%+v)", value, value)
}

func asBool(value interface{}) (bool, error) {
	result, ok := value.(bool)
	if ok {
		return result, nil
	}

	return false, fmt.Errorf("expecting a boolean, found a %T (%+v)", value, value)
}

func asListOfStrings(value interface{}) ([]string, error) {
	if v, ok := value.([]string); ok {
		return v, nil
	}

	elements, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expecting a list of strings, found a %T (%+v)", value, value)
	}

	listOfStrings := []string{}

	for _, element := range elements {
		str, ok := element.(string)
		if !ok {
			return nil, fmt.Errorf("expecting a list of strings, found an invalid element: type %T (%+v)", element, element)
		}
		listOfStrings = append(listOfStrings, str)
	}

	return listOfStrings, nil
}
