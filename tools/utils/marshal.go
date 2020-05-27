// Author: yann
// Date: 2020/5/23 11:32 上午
// Desc:

package utils

import "encoding/json"

func JsonMarshal(i interface{}) string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func JsonUnMarshal(jsonStr []byte, out interface{}) error {
	return json.Unmarshal(jsonStr, out)
}
