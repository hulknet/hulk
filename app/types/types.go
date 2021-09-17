package types

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math"
	"math/big"
	"math/bits"
	rd "math/rand"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ErrDecodeByte32 = "decode of ByteArray failed"
	ErrSizeByte32   = "size of ByteArray is invalid"
	ErrSizeSign     = "size of Sign is invalid"

	ErrGetToken    = "failed to get Token from request header"
	ErrDecodeToken = "decode of Token failed"
	ErrGetAddr     = "failed to get Address  from request header"
	ErrGetID       = "failed to get ID  from request header"
	ErrDecodeAddr  = "decode of Address failed"
	ErrGetSign     = "failed to get Signature  from request header"
	ErrDecodeSign  = "decode of Signature failed"
	ErrDecodePart  = "decode of  Partition failed"
)

type Token [32]byte
type Permission []byte
type Peer struct {
	PK    PK
	Token Token
}

func (p Peer) Equal(other Peer) bool {
	return p.PK == other.PK
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

type ID [32]byte

func (i ID) IsEmpty() bool {
	var empty ID
	return i == empty
}

func (i ID) Addr() Addr {
	return Addr(binary.BigEndian.Uint64(i[:8]))
}

type PK [32]byte

func (p PK) ID() ID {
	var id ID
	copy(id[:], p[:])
	return id
}

type Sign [65]byte

func (s Sign) PK() PK {
	var pk PK
	copy(pk[:], s[32:64])
	return pk
}

func ToHex(s [32]byte) string {
	return hex.EncodeToString(s[:])
}

func FromHex(s string) ([32]byte, error) {
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

func SignFromHex(s string) (Sign, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return Sign{}, errors.New(ErrDecodeSign)
	}

	var sign Sign
	bitLen := copy(sign[:], data[:65])
	if bitLen != 65 {
		return Sign{}, errors.New(ErrSizeSign)
	}

	return sign, nil
}

type Addr uint64

func (a Addr) Cpl(target Addr) int {
	return bits.LeadingZeros64(uint64(a ^ target))
}

// Normalize between 0 and 1
func (a Addr) Normalize(bitSizePrefix uint8) float64 {
	val := clearBitPrefix(uint64(a), bitSizePrefix)
	max := clearBitPrefix(math.MaxUint64, bitSizePrefix)
	return float64(val) / float64(max)
}

func clearBitPrefix(val uint64, bitSizePrefix uint8) uint64 {
	return (val << bitSizePrefix) >> bitSizePrefix
}

// pkcs1PublicKey reflects the ASN.1 structure of a PKCS #1 public key.
type pkcs1PublicKey struct {
	N *big.Int
	E int
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func DecodeDefaultPublicKey() (*ecdsa.PublicKey, error) {
	path, err := filepath.Abs("./app/pubKey.pem")
	if err != nil {
		return nil, err
	}
	pemEncodedPub, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	blockPub, _ := pem.Decode(pemEncodedPub)
	encodedPub := blockPub.Bytes
	var pki publicKeyInfo
	if rest, err := asn1.Unmarshal(encodedPub, &pki); err != nil {
		if _, err := asn1.Unmarshal(encodedPub, &pkcs1PublicKey{}); err == nil {
			return nil, errors.New("x509: failed to parse public key (use ParsePKCS1PublicKey instead for this key format)")
		}
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after ASN.1 of public-key")
	}

	publicKey, err := crypto.UnmarshalPubkey(pki.PublicKey.RightAlign())
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

func DecodeDefaultPrivateKey() (*ecdsa.PrivateKey, error) {
	path, err := filepath.Abs("./app/privateKey.pem")
	if err != nil {
		return nil, err
	}
	pemEncoded, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemEncoded)
	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(block.Bytes, &privKey); err != nil {
		return nil, errors.New("x509: failed to parse EC private key: " + err.Error())
	}
	privateKey, err := crypto.ToECDSA(privKey.PrivateKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
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
