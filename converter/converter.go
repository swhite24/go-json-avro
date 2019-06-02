package converter

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Convert converts a standard json structure into an encodable structure suitable
// for the provided avro schema.
func Convert(datum map[string]interface{}, schema string) (map[string]interface{}, error) {
	var err error
	var parsedSchema, converted map[string]interface{}
	var out interface{}

	if err = json.Unmarshal([]byte(schema), &parsedSchema); err != nil {
		return nil, err
	}

	if out, err = derive(datum, parsedSchema); err != nil {
		return nil, err
	}

	converted, _ = out.(map[string]interface{})
	return converted, nil
}

// derive determines the appropriate structure of the provided value for the given type
func derive(datum, typ interface{}) (interface{}, error) {
	switch t := typ.(type) {
	// Simple case. typ contains the primitive type name to set.
	case string:
		if t == "null" && datum != nil {
			return nil, errors.New("value cannot match null type")
		}
		return datum, nil

	// Assume record type for now
	case map[string]interface{}:
		var data map[string]interface{}
		var fields []interface{}
		var ok bool
		var typ string
		var out = map[string]interface{}{}

		// Short circuit to try and save non record types
		if typ, ok = t["type"].(string); ok && typ != "record" {
			return derive(datum, typ)
		}

		if data, ok = datum.(map[string]interface{}); !ok {
			return nil, errors.New("invalid type provided for record field")
		}

		if fields, ok = t["fields"].([]interface{}); !ok {
			return nil, errors.New("invalid format for fields")
		}

		// Handle each field, constructing map
		for _, f := range fields {
			var field map[string]interface{}
			var name string
			var ok bool

			if field, ok = f.(map[string]interface{}); !ok {
				continue
			}

			if name, ok = field["name"].(string); !ok {
				continue
			}

			// Recursively capture value
			if val, err := derive(data[name], field["type"]); err == nil {
				out[name] = val
			}
		}
		return out, nil

	// Union type
	case []interface{}:
		// Nothing special if omitted / nil
		if datum == nil {
			return nil, nil
		}

		out := map[string]interface{}{}

		// Check each potential type
		for _, potential := range t {
			// Check if provided value is determined to be valid for potential type
			if val, err := derive(datum, potential); err == nil {
				// Value is valid, get correct name format and return
				if name, err := getTypeName(potential); err == nil {
					out[name] = val
					return out, nil
				}
			}
		}
	}
	return nil, nil
}

// getTypeName determines the name a value should receive when part of a union type
func getTypeName(t interface{}) (string, error) {
	switch typ := t.(type) {
	// Primitive / known type, return t
	case string:
		return typ, nil
	// Record type, construct name from record name and optionally namespace
	case map[string]interface{}:
		var pkg, name string
		var ok bool

		if name, ok = typ["name"].(string); !ok {
			return "", errors.New("invalid name provided")
		}

		if typ["namespace"] != nil {
			if pkg, ok = typ["namespace"].(string); !ok {
				return "", errors.New("invalid namespace provided")
			}

			return fmt.Sprintf("%s.%s", pkg, name), nil
		}
		return name, nil
	}
	return "", errors.New("unrecognized type")
}
