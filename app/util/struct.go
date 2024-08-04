package util

import (
	"encoding/json"
	"fmt"
	"html"
)

func StructToJsonString(data interface{}) string {

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return html.UnescapeString(string(b))
}

func GetStringByMap(dataMap map[string]interface{}, key string) (res string) {

	data, ok := dataMap[key]
	if ok {
		if res, ok = data.(string); ok {
			return res
		}
	}

	return ""
}
