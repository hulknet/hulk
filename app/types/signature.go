package types

import (
	"encoding/hex"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

type Sign [65]byte

func (s Sign) SignWithPK() (sign [65]byte) {
	copy(sign[:], s[:])
	return
}

func (s Sign) Sign() (sign [64]byte) {
	copy(sign[:], s[:64])
	return
}

func (s Sign) PK(hash []byte) (pk PK, err error) {
	pkSource, err := secp256k1.RecoverPubkey(hash, s[:])
	if err != nil {
		return
	}
	pk.Write(pkSource)

	return
}

func (s Sign) CheckSignature(msg []byte) (bool, error) {
	msgHash := sha3.Sum256(msg)
	return s.CheckHashSignature(msgHash[:])
}

func (s Sign) CheckHashSignature(hash []byte) (bool, error) {
	pk, err := s.PK(hash)
	if err != nil {
		return false, err
	}

	return secp256k1.VerifySignature(pk[:], hash[:], s[:64]), nil
}

func CheckSignature(msg []byte, sign []byte) (bool, error) {
	msgHash := sha3.Sum256(msg)
	pk, err := secp256k1.RecoverPubkey(msgHash[:], sign[:])
	if err != nil {
		return false, err
	}

	return secp256k1.VerifySignature(pk, msgHash[:], sign[:64]), nil
}

func CheckStringSignature(msg []byte, strSign string) (bool, error) {
	sign, err := SignFromHex(strSign)
	if err != nil {
		return false, err
	}

	msgHash := sha3.Sum256(msg)
	pk, err := secp256k1.RecoverPubkey(msgHash[:], sign[:])
	if err != nil {
		return false, err
	}

	return secp256k1.VerifySignature(pk, msgHash[:], sign[:64]), nil
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
