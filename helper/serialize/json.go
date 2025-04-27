package serialize

import "encoding/json"

func JsonSerialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func JsonDeserialize(b []byte, obj interface{}) error {
	return json.Unmarshal(b, &obj)
}
