package types

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	rd "math/rand"
	"time"
)

const (
	ErrDecodeByte32 = "decode of Byte32 failed"
	ErrDecodeByte8  = "decode of Byte8 failed"
	ErrDecodeByte   = "decode of ByteArray failed"
	ErrSizeByte32   = "size of Byte32 is invalid"
	ErrSizeByte8    = "size of Byte8 is invalid"
	ErrSizeByte     = "size of ByteArray is invalid"
	ErrSizeSign     = "size of Sign is invalid"

	ErrGetToken    = "failed to get Token from request header"
	ErrDecodeToken = "decode of Token failed"
	ErrGetAddr     = "failed to get Address  from request header"
	ErrGetID       = "failed to get ID  from request header"
	ErrGetIDPrefix = "failed to get ID  from request header"
	ErrDecodeAddr  = "decode of Address failed"
	ErrGetSign     = "failed to get Signature  from request header"
	ErrDecodeSign  = "decode of Signature failed"
	ErrDecodePart  = "decode of  Partition failed"
)

type Block struct {
	ID      ID
	PID     ID
	BitSize []uint8
}

type Token [32]byte
type Permission []byte
type Peer struct {
	PK    PK
	Token Token
}

func (p Peer) Equal(other Peer) bool {
	return p.PK == other.PK
}

type ID [32]byte

func (i ID) IsEmpty() bool {
	var empty ID
	return i == empty
}

func (i ID) Addr() Addr {
	return Addr(binary.BigEndian.Uint64(i[:8]))
}

func (i ID) Prefix() (prefix IDPrefix) {
	copy(prefix[:], i[:8])
	return
}

type IDPrefix [8]byte

func (ip IDPrefix) Bytes() []byte {
	return ip[:]
}

func (ip IDPrefix) Addr() Addr {
	return Addr(binary.BigEndian.Uint64(ip[:]))
}

type PK [32]byte

func (p PK) ID() (id ID) {
	copy(id[:], p[:])
	return id
}

func (p PK) Bytes() []byte {
	return p[:]
}

func (p PK) Prefix() (prefix IDPrefix) {
	copy(prefix[:], p[:])
	return
}

func (p PK) Write(data []byte) {
	copy(p[:], data[:])
}

func ToHex(s [32]byte) string {
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

func IDFromHex(s string) ([32]byte, error) {
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

func IDPrefixFromHex(s string) (idPrefix IDPrefix, err error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return idPrefix, errors.New(ErrDecodeByte8)
	}

	bitLen := copy(idPrefix[:], data[:8])
	if bitLen != 8 {
		err = errors.New(ErrSizeByte8)
	}

	return
}

func Random() int {
	rd.Seed(time.Now().UnixNano())
	return rd.Int()
}

func GenerateSHA() [32]byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(Random()))
	return sha256.Sum256(bs)
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
