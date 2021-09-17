package types

import (
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

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
