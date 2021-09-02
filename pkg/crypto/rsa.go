package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/pkg/errors"
)

type RSAPublicKey rsa.PublicKey

func (pk *RSAPublicKey) Verify(data, sig []byte) (bool, error) {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15((*rsa.PublicKey)(pk), crypto.SHA256, hashed[:], sig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pk *RSAPublicKey) Marshal() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pk)
}

func (pk *RSAPublicKey) ID() (ID, error) {
	return rsaPublicKeyToID((*rsa.PublicKey)(pk))
}

func unmarshalRsaPublicKey(b []byte) (*RSAPublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	pk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("can't unmarshal rsa public key")
	}
	if pk.N.BitLen() < 256 {
		return nil, errors.New("rsa public key is not 256bit")
	}
	return (*RSAPublicKey)(pk), nil
}

func rsaPublicKeyToID(k *rsa.PublicKey) (ID, error) {
	data, err := x509.MarshalPKIXPublicKey(k)
	if err != nil {
		return ID{}, errors.Wrap(err, "failed to marshal rsa public key")
	}
	return sha256.Sum256(data), nil
}

type RSAPrivateKey rsa.PrivateKey

func (p *RSAPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, (*rsa.PrivateKey)(p), crypto.SHA256, h[:])
}

func (p *RSAPrivateKey) Marshal() ([]byte, error) {
	return x509.MarshalPKCS1PrivateKey((*rsa.PrivateKey)(p)), nil
}

func (p *RSAPrivateKey) Public() PublicKey {
	pk := (*rsa.PrivateKey)(p)
	return (*RSAPublicKey)(&pk.PublicKey)
}

func unmarshalRsaPrivateKey(b []byte) (*RSAPrivateKey, error) {
	pk, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return (*RSAPrivateKey)(pk), nil
}
