package bdd

type Encodable interface {
	Encode() ([]byte, []byte)
}
type Decodable interface {
	Decode([]byte) error
}
