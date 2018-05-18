package bdd

type Serializable interface {
	Encodable
	Decodable
}

type Encodable interface {
	Encode() []byte
}
type Decodable interface {
	Decode([]byte) error
}
