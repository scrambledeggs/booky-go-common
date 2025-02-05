package converters

import "encoding/json"

func StructToInterface[T any](data T) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(data)
	json.Unmarshal(inrec, &inInterface)

	return inInterface
}
