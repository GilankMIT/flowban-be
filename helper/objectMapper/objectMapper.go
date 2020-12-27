package objectMapper

import "encoding/json"

//StructToMap convert struct to map
func StructToMap(fromStruct interface{}) (map[string]interface{}, error) {

	var inInterface map[string]interface{}

	inrec, err := json.Marshal(&fromStruct)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}

	return inInterface, nil
}
