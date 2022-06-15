package format

import (
	"encoding/json"
	"fmt"
)

// Print object details to json format
func PrintObject(value interface{}) {
	data, _ := json.MarshalIndent(value, "", "    ")
	fmt.Println(string(data))
}
