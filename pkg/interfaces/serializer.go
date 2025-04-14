package interfaces

type Serializer[T any] interface {
	// Serialize the node to a byte array
	Serialize(param T) ([]byte, error)
	// Deserialize the byte array to a node
	Deserialize(data []byte) (T, error)
}
