package core

import (
	"fmt"
	"reflect"
)

// Convert and apply args to arbitrary function fn
func Cumin(fn interface{}, args []interface{}) ([]interface{}, error) {
	reciever := reflect.TypeOf(fn)
	var ret []interface{}

	if reciever.Kind() != reflect.Func {
		return ret, fmt.Errorf("Handler is not a function!")
	}

	if reciever.NumIn() != len(args) {
		return ret, fmt.Errorf("Cumin: expected %s args for function %s, got %s", reciever.NumIn(), reciever, len(args))
	}

	// Iterate over the params listed in the method and try their casts
	values := make([]reflect.Value, len(args))

	for i := 0; i < reciever.NumIn(); i++ {
		param := reciever.In(i)
		arg := reflect.ValueOf(args[i])

		if param == arg.Type() {
			values[i] = arg
		} else if arg.Type().ConvertibleTo(param) {
			values[i] = arg.Convert(param)
		} else {
			return ret, fmt.Errorf("Cumin: expected %s for arg[%d] in (%s), got %s.", param, i, reciever, arg.Type())
		}
	}

	// Perform the call
	result := reflect.ValueOf(fn).Call(values)
	for _, x := range result {
		ret = append(ret, x.Interface())
	}

	// Catch any exceptions this produces and pass them to the function that sent them, or some kind of handler

	return ret, nil
}

/*
   There are four cases softCumin can hit.
       - Initial case: positional array of heterogenous elements (starting condition)
       - Lists: array of homogenous elements checked against one type
       - Dictionaries: dictionary of key/value pairs. Check key, value, and recursive
       - Primitive: end case

   The type of the next case is determined by reflecting the type of the expected element at that position.
   Current version is not recursive, so dictionaries are not handled.

Version from the test:
    [int [str] map[a:str b:int]] against [1 [Hey There] map[b:1 a:alpha]]

Version from swift:
    [int [str] map[age:int name:str]]
    [int [str] {name: str, age: int}] against [1 [Hey There] map[name:Billiam]]
*/

// Checks the types of the provided positional arguments and the receiver.
func softCumin(types []interface{}, args []interface{}) error {
	// fmt.Printf("SOFTCUMIN: %v against %v\n", types, args)

	// Description to print on every failure
	description := fmt.Sprintf("Types: %v, Arguments: %v", types, args)

	if len(types) != len(args) {
		return fmt.Errorf("Cumin: Invalid number of arguments, expected %d, receieved %s. %v", len(types), len(args), description)
	}

	for i, x := range args {
		argument := reflect.ValueOf(x)
		expected := types[i]

		// Debug("Expected: %v expected type: %v, Argument: %v", expected, reflect.TypeOf(expected), x)
		//fmt.Printf("Expected: %v expected type: %v, Argument: %v\n", expected, reflect.TypeOf(expected), x)

		// If the expected type is a string, we're looking for a primitive
		if s, ok := expected.(string); ok {
			if e := primitiveCheck(s, argument.Kind()); e != nil {
				return e
			}

		} else if nestedSlice, ok := expected.([]interface{}); ok {
			if len(nestedSlice) != 1 {
				return fmt.Errorf("Cumin: array expected at position #%d is not homogenous. %s", i, expected)
			}

			if argumentList, ok := x.([]interface{}); !ok {
				return fmt.Errorf("Cant read interface list %v at position %d", x, i)
			} else {
				for _, v := range argumentList {
					if e := primitiveCheck(nestedSlice[0].(string), reflect.TypeOf(v).Kind()); e != nil {
						return e
					}
				}
			}

		} else if nestedMap, ok := expected.(map[string]interface{}); ok {

			if argumentMap, ok := x.(map[string]interface{}); !ok {
				return fmt.Errorf("Cumin: expected dictionary at position %d, got %v", i, reflect.TypeOf(x))
			} else {
				if e := mapCheck(nestedMap, argumentMap); e != nil {
					return e
				}
			}
		} else {
			return fmt.Errorf("Cumin: couldnt find primitive, list, or dictionary at #%d. %v", i, description)
		}
	}

	return nil
}

// Return an error if the argument is not of the expected type OR the expected type is not a primitive
func primitiveCheck(expected string, argument reflect.Kind) error {
	if argument == reflect.Bool && expected == "bool" ||
		argument == reflect.String && expected == "str" ||
		argument == reflect.Float64 && (expected == "float" || expected == "int") ||
		argument == reflect.Int && (expected == "float" || expected == "int") {
		return nil
	}

	return fmt.Errorf("Cumin: expecting primitive %s, got %s", expected, argument)
}

// Recursively check an object. Return nil if the object matches the expected types
func mapCheck(expected map[string]interface{}, argument map[string]interface{}) error {
	if len(expected) != len(argument) {
		return fmt.Errorf("Cumin: object invalid number of keys, expected %d, receieved %s", len(expected), len(argument))
	}

	// TODO: nested collections and objects
	for k, v := range argument {
		if e := primitiveCheck(expected[k].(string), reflect.ValueOf(v).Kind()); e != nil {
			return e
		}
	}

	return nil
}
