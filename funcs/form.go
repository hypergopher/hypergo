package funcs

import "fmt"

func attrsToMap(data map[string]any, specialAttrs map[string]string, attrs ...any) (map[string]any, error) {
	attributes := make(map[string]string)

	// Always add the following, if they aren't already there
	if _, found := data["Hyperscript"]; !found {
		data["Hyperscript"] = ""
	}

	if _, found := data["Class"]; !found {
		data["Class"] = ""
	}

	// Special attrs that are always added for all fields
	if _, found := specialAttrs["class"]; !found {
		specialAttrs["class"] = "Class"
		specialAttrs["_"] = "Hyperscript"
		specialAttrs["hyperscript"] = "Hyperscript"
	}

	for i := 0; i < len(attrs); i += 2 {
		key, ok := attrs[i].(string)
		if !ok {
			return nil, fmt.Errorf("[InputAttrs] attribute key at position %d is not a string", i)
		}
		value, ok := attrs[i+1].(string)
		if !ok {
			// TODO: Add support for other types? Or, just ignore?
			// For now, we just ignore non-string values.
			value = ""
		}

		// Check if the attribute is a special attribute
		if sk, found := specialAttrs[key]; found {
			data[sk] = value
		} else {
			attributes[key] = value
		}
	}

	data["Attributes"] = attributes

	return data, nil
}

// InputAttrs prepares the data for rendering an input field in a separate template.
// nameID is used for both the name and ID attributes for simplicity.
// errors is a slice of strings representing validation messages for this input.
// attrs are additional attributes provided as key/value pairs.
// Use a "label" attribute to set the label text, and a "type" attribute to set the input type.
// Use a "hint" attribute to set the hint text.
// Use an "error" attribute to set an error message.
// Use a "togglePassword" attribute to "true" to add a toggle button for password visibility (only for password inputs).
func InputAttrs(nameID string, attrs ...any) (map[string]any, error) {
	if len(attrs)%2 != 0 {
		return nil, fmt.Errorf("InputAttrs expects attributes as key/value pairs, received odd number of arguments")
	}

	data := map[string]any{
		"NameID": nameID,
		"Error":  "",
		"Hint":   "",
		"Label":  "",
		"Type":   "text",
	}

	// attributes := make(map[string]string)
	specialAttrs := map[string]string{
		"error": "Error",
		"hint":  "Hint",
		"label": "Label",
		"type":  "Type",
	}

	data, err := attrsToMap(data, specialAttrs, attrs...)
	if err != nil {
		return nil, err
	}

	return data, nil
}
