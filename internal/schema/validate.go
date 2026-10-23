package schema

import (
	"fmt"
	"time"
)

// ValidateData checks a record's data against its object's field definitions:
// required presence (falling back to a field's default), and a basic type
// check per field type. It returns the normalized data (defaults applied,
// file/readonly fields stripped since those are never client-supplied).
//
// Not validated here: relation fields are checked for type only (a number),
// not that the referenced row exists — add that if orphaned relations become
// a real problem. current_stock-style calculated fields are computed later
// by reporting, not derived at write time.
func ValidateData(fields []FieldDefinition, data map[string]any) (map[string]any, error) {
	result := map[string]any{}

	for _, f := range fields {
		if f.Type == "file" {
			continue
		}
		meta := f.ParseMetadata()
		if meta.Readonly || meta.Calculated {
			continue
		}

		value, present := data[f.Name]
		if !present || value == nil || value == "" {
			if meta.Default != "" {
				result[f.Name] = meta.Default
				continue
			}
			if f.Required {
				return nil, fmt.Errorf("%s is required", f.Name)
			}
			continue
		}

		normalized, err := checkType(f, meta, value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name, err)
		}
		result[f.Name] = normalized
	}

	return result, nil
}

func checkType(f FieldDefinition, meta FieldMetadata, value any) (any, error) {
	switch f.Type {
	case "text":
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("must be a string")
		}
		return s, nil
	case "number", "relation":
		n, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("must be a number")
		}
		return n, nil
	case "date":
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("must be a date string (YYYY-MM-DD)")
		}
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return nil, fmt.Errorf("must be a valid date (YYYY-MM-DD)")
		}
		return s, nil
	case "time":
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("must be a time string (HH:MM)")
		}
		if _, err := time.Parse("15:04", s); err != nil {
			return nil, fmt.Errorf("must be a valid time (HH:MM)")
		}
		return s, nil
	case "datetime":
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("must be an RFC3339 datetime string")
		}
		if _, err := time.Parse(time.RFC3339, s); err != nil {
			return nil, fmt.Errorf("must be a valid RFC3339 datetime")
		}
		return s, nil
	case "select":
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("must be a string")
		}
		if len(meta.Options) > 0 && !contains(meta.Options, s) {
			return nil, fmt.Errorf("must be one of %v", meta.Options)
		}
		return s, nil
	case "multiselect":
		items, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("must be an array of strings")
		}
		values := make([]string, 0, len(items))
		for _, item := range items {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("must be an array of strings")
			}
			if len(meta.Options) > 0 && !contains(meta.Options, s) {
				return nil, fmt.Errorf("%q must be one of %v", s, meta.Options)
			}
			values = append(values, s)
		}
		return values, nil
	default:
		return value, nil
	}
}

func contains(options []string, want string) bool {
	for _, o := range options {
		if o == want {
			return true
		}
	}
	return false
}
