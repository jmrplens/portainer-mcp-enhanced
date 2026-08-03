// Package toolgen provides tool generation utilities for building MCP tool
// definitions from YAML specifications. It handles parameter extraction,
// type conversion, and tool definition construction.
package toolgen

import (
	"fmt"
	"math"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// ParameterParser provides methods to safely extract parameters from request arguments
type ParameterParser struct {
	args map[string]any
}

// NewParameterParser creates a new parameter parser for the given request
func NewParameterParser(request mcp.CallToolRequest) *ParameterParser {
	return &ParameterParser{
		args: request.GetArguments(),
	}
}

// GetString extracts a string parameter from the request
func (p *ParameterParser) GetString(name string, required bool) (string, error) {
	value, ok := p.args[name]
	if !ok || value == nil {
		if required {
			return "", fmt.Errorf("%s is required", name)
		}
		return "", nil
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", name)
	}

	return strValue, nil
}

// GetNumber extracts a number parameter from the request. As a fallback for
// clients that don't declare (and therefore can't type-check) this parameter
// — e.g. a grouped meta-tool whose schema only advertised an "action" enum
// before per-action schemas were merged in — a numeric string like "55" is
// also accepted.
func (p *ParameterParser) GetNumber(name string, required bool) (float64, error) {
	value, ok := p.args[name]
	if !ok || value == nil {
		if required {
			return 0, fmt.Errorf("%s is required", name)
		}
		return 0, nil
	}

	numValue, ok := coerceFloat64(value)
	if !ok {
		return 0, fmt.Errorf("%s must be a number", name)
	}

	return numValue, nil
}

// coerceFloat64 extracts a float64 from a decoded JSON value, additionally
// accepting a numeric string (e.g. "55") as a fallback for callers that sent
// an untyped parameter as text.
func coerceFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// GetInt extracts an integer parameter from the request
func (p *ParameterParser) GetInt(name string, required bool) (int, error) {
	num, err := p.GetNumber(name, required)
	if err != nil {
		return 0, err
	}

	if num < math.MinInt || num > math.MaxInt || math.Trunc(num) != num {
		return 0, fmt.Errorf("%s must be a valid integer", name)
	}

	return int(num), nil
}

// GetBoolean extracts a boolean parameter from the request. As a fallback for
// clients that can't type-check this parameter (see GetNumber), a boolean
// string like "true" or "false" is also accepted.
func (p *ParameterParser) GetBoolean(name string, required bool) (bool, error) {
	value, ok := p.args[name]
	if !ok || value == nil {
		if required {
			return false, fmt.Errorf("%s is required", name)
		}
		return false, nil
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("%s must be a boolean", name)
		}
		return b, nil
	default:
		return false, fmt.Errorf("%s must be a boolean", name)
	}
}

// GetArrayOfIntegers extracts an array of numbers parameter from the request
func (p *ParameterParser) GetArrayOfIntegers(name string, required bool) ([]int, error) {
	value, ok := p.args[name]
	if !ok || value == nil {
		if required {
			return nil, fmt.Errorf("%s is required", name)
		}
		return []int{}, nil
	}

	arrayValue, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be an array", name)
	}

	return parseArrayOfIntegers(arrayValue)
}

// GetArrayOfObjects extracts an array of objects parameter from the request
func (p *ParameterParser) GetArrayOfObjects(name string, required bool) ([]any, error) {
	value, ok := p.args[name]
	if !ok || value == nil {
		if required {
			return nil, fmt.Errorf("%s is required", name)
		}
		return []any{}, nil
	}

	arrayValue, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be an array", name)
	}

	return arrayValue, nil
}

// parseArrayOfIntegers converts a slice of any type to a slice of integers.
// Each element may be a JSON number or (as a fallback for untyped callers,
// see GetNumber) a numeric string. Returns an error if any value cannot be
// parsed as an integer.
//
// Example:
//
//	ids, err := parseArrayOfIntegers([]any{1, 2, 3})
//	// ids = []int{1, 2, 3}
func parseArrayOfIntegers(array []any) ([]int, error) {
	result := make([]int, 0, len(array))

	for _, item := range array {
		idFloat, ok := coerceFloat64(item)
		if !ok {
			return nil, fmt.Errorf("failed to parse '%v' as integer", item)
		}
		if idFloat < math.MinInt64 || idFloat > math.MaxInt64 || math.Trunc(idFloat) != idFloat {
			return nil, fmt.Errorf("failed to parse '%v' as integer", item)
		}
		result = append(result, int(idFloat))
	}

	return result, nil
}
