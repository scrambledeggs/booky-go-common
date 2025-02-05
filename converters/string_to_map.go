package converters

import "encoding/json"

func StringToMap(rawData string) map[string]any {
	var mapData map[string]any

	if err := json.Unmarshal([]byte(rawData), &mapData); err != nil {
		return map[string]any{}
	}

	return mapData
}
