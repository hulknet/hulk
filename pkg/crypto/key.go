package crypto

import (
	"bytes"
)

type Key interface {
	Marshal() ([]byte, error)
}

type PublicKey interface {
	Key
	ID() (ID, error)
	Verify(data, signature []byte) (bool, error)
}

type PrivateKey interface {
	Key
	Public() PublicKey
	Sign(data []byte) ([]byte, error)
}

func UnmarshalPublicKey(data []byte) (PublicKey, error) {
	return unmarshalRsaPublicKey(data)
}

func UnmarshalPrivateKey(data []byte) (PrivateKey, error) {
	return unmarshalRsaPrivateKey(data)
}

func KeyEqual(k1, k2 Key) (bool, error) {
	b1, err := k1.Marshal()

	if err != nil {
		return false, err
	}
	b2, err := k2.Marshal()
	if err != nil {
		return false, err
	}
	return bytes.Equal(b1, b2), nil
}
