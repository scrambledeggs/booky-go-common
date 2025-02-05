package converters

import (
	"encoding/json"
)

func MapToString(rawData map[string]any) string {
	str, err := json.Marshal(rawData)

	if err != nil {
		return ""
	}

	return string(str)
}
