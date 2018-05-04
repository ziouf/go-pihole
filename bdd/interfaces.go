package bdd


type Encodable interface {
	Encode() ([]byte, error)
	StampEncoded() []byte
}
type Decodable interface {
	Decode([]byte) error
}

