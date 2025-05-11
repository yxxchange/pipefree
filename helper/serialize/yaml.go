package serialize

import "gopkg.in/yaml.v3"

func YamlSerialize(obj interface{}) ([]byte, error) {
	return yaml.Marshal(obj)
}

func YamlDeserialize(b []byte, obj interface{}) error {
	return yaml.Unmarshal(b, obj)
}
