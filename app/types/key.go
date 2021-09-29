package types

import (
	"crypto/ecdsa"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
)

type PubKey [65]byte

func (p PubKey) ID256() (id ID256) {
	copy(id[:], p[:32])
	return id
}

func (p PubKey) Bytes() []byte {
	return p[:]
}

func (p PubKey) ID() (id ID64) {
	copy(id[:], p[:8])
	return
}

func (p *PubKey) Write(data []byte) {
	copy(p[:], data[:])
}

type ECKey struct {
	privateKey *ecdsa.PrivateKey
}

func (k *ECKey) ECPrivateKey() *ecdsa.PrivateKey {
	return k.privateKey
}

func (k *ECKey) Bytes() []byte {
	return crypto.FromECDSA(k.privateKey)
}

func (k *ECKey) Pub() (pub PubKey) {
	pub.Write(crypto.FromECDSAPub(k.privateKey.Public().(*ecdsa.PublicKey)))
	return
}

func HexToECKey(hex string) (*ECKey, error) {
	privateKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, err
	}
	return &ECKey{
		privateKey: privateKey,
	}, nil
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
