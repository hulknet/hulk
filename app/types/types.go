package types

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

const (
	ErrDecodeByte32  = "decode of Byte32 failed"
	ErrDecodeByte8   = "decode of Byte8 failed"
	ErrDecodeByte    = "decode of ByteArray failed"
	ErrDecodeTime    = "decode of Time failed"
	ErrDecodeSign520 = "decode of Signature failed"
	ErrDecodePart    = "decode of Partition failed"
	ErrDecodeToken   = "decode of Token failed"

	ErrInvalidTime = "object Time is invalid"
	ErrSizeByte32  = "size of Byte32 is invalid"
	ErrSizeByte8   = "size of Byte8 is invalid"
	ErrSizeByte    = "size of ByteArray is invalid"
	ErrSizeSign520 = "size of Sign520 is invalid"

	ErrGetToken = "failed to get Token from request header"
	ErrGetTime  = "failed to get Time from request header"
	ErrGetID    = "failed to get ID from request header"
	ErrGetSign  = "failed to get Signature from request header"
)

type Token [32]byte
type Permission []byte
type Peer struct {
	Pub   PubKey
	Token Token
}

func (p Peer) Equal(other Peer) bool {
	return p.Pub == other.Pub
}

type ID256 [32]byte

func (i ID256) IsEmpty() bool {
	var empty ID256
	return i == empty
}

func (i ID256) Uint64() uint64 {
	return binary.BigEndian.Uint64(i[:8])
}

func (i ID256) Bytes() []byte {
	return i[:]
}

func (i ID256) ID64() (prefix ID64) {
	copy(prefix[:], i[:8])
	return
}

type ID64 [8]byte

func (i ID64) Bytes() []byte {
	return i[:]
}

func (i ID64) Uint64() uint64 {
	return binary.BigEndian.Uint64(i[:])
}

func ID256ToHex(s [32]byte) string {
	return hex.EncodeToString(s[:])
}

func FromHex(s string, byteLen int) ([]byte, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}, errors.New(ErrDecodeByte)
	}

	if byteLen != 0 && byteLen != len(data) {
		return []byte{}, errors.New(ErrSizeByte)
	}

	return data, nil
}

func ID256FromHex(s string) ([32]byte, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, errors.New(ErrDecodeByte32)
	}

	var id [32]byte
	bitLen := copy(id[:], data[:32])
	if bitLen != 32 {
		return [32]byte{}, errors.New(ErrSizeByte32)
	}

	return id, nil
}

func ID64FromHex(s string) (id ID64, err error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return id, errors.New(ErrDecodeByte8)
	}

	bitLen := copy(id[:], data[:8])
	if bitLen != 8 {
		err = errors.New(ErrSizeByte8)
	}

	return
}

type Bitmap256 [8]byte

func (b *Bitmap256) IsSet(i byte) bool { i -= 1; return b[i/8]&(1<<uint(7-i%8)) != 0 }
func (b *Bitmap256) Set(i byte)        { i -= 1; b[i/8] |= 1 << uint(7-i%8) }
func (b *Bitmap256) Clear(i byte)      { i -= 1; b[i/8] &^= 1 << uint(7-i%8) }

func (b *Bitmap256) Sets(xs ...byte) {
	for _, x := range xs {
		b.Set(x)
	}
}
