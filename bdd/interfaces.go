package bdd

type Encodable interface {
	Encode() []byte
}
type Decodable interface {
	Decode([]byte) error
}
