package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	ErrDecodeID = "decode of ID failed"
	ErrSizeID   = "ID size is invalid"

	IDLogLen = 80
)

type ID [32]byte

func (i ID) IsEmpty() bool {
	var empty ID
	return i == empty
}

// Slice convert [32]byte to []byte
func (i ID) Slice() []byte {
	return i[:]
}

// SliceL copy first numbers (length) of bits
func (i ID) SliceL(length int) []byte {
	fullBites := length / 8
	tailBits := 8 - uint8(length%8)
	key := make([]byte, 0)
	if fullBites > 0 {
		key = append(key, i[:fullBites]...)
	}
	if tailBits > 0 {
		shiftedBite := i[fullBites] >> tailBits
		key = append(key, shiftedBite)
	}

	return key
}

func (i ID) Hex() string {
	return hex.EncodeToString(i[:])
}

func (i ID) HexL(length int) string {
	return hex.EncodeToString(i.SliceL(length))
}

func (i ID) Replica(num int) ID {
	if num <= 0 {
		return i
	}
	return ID(sha256.Sum256(i[:])).Replica(num - 1)
}

func (i ID) WithSalt(salt []byte) ID {
	return sha256.Sum256(append(i[:], salt...))
}

func FromHex(s string) (ID, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return ID{}, errors.New(ErrDecodeID)
	}

	var id ID
	bitLen := copy(id[:], data[:32])
	if bitLen != 32 {
		return ID{}, errors.New(ErrSizeID)
	}

	return id, nil
}

type Signature [64]byte

func (s Signature) Hex() string {
	return hex.EncodeToString(s[:])
}

func SignatureFromHex(s string) (Signature, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return Signature{}, errors.New(ErrDecodeID)
	}

	var sign Signature
	bitLen := copy(sign[:], data[:64])
	if bitLen != 64 {
		return Signature{}, errors.New(ErrSizeID)
	}

	return sign, nil
}
