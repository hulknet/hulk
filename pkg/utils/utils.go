package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"io"
	"math/big"
	"math/bits"
	rd "math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	cr "github.com/kotfalya/hulk/pkg/crypto"
)

func TryClose(l *log.Entry, closer ...io.Closer) {
	for _, c := range closer {
		if err := c.Close(); err != nil {
			l.WithError(err).Error("error on close")
		}
	}
}

func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{tlsCert},
		NextProtos:         []string{"hulk-net"},
	}
}

func Random() int {
	rd.Seed(time.Now().UnixNano())
	return rd.Int()
}

func GenerateSHA() cr.ID {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(Random()))
	return sha256.Sum256(bs)
}

func GenerateSHAFrom(source string) cr.ID {
	bs := []byte(source)
	return sha256.Sum256(bs)
}

func Cpl(p1, p2 []byte) int {
	k3 := XOR(p1, p2)
	return ZeroPrefixLen(k3)
}

func ZeroPrefixLen(id []byte) int {
	for i, b := range id {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}

	return len(id) * 8
}

func XOR(a, b []byte) []byte {
	c := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func OneCount(p1, p2 []byte) int {
	p3 := XOR(p1, p2)
	var r int
	for i := 0; i < len(p3); i++ {
		r += bits.OnesCount8(p3[i])
	}
	return r
}

func Distance(p1, p2 []byte) *big.Int {
	k3 := XOR(p1, p2)
	dist := big.NewInt(0).SetBytes(k3)

	return dist
}

func IdToAddr(shift uint8, id cr.ID) uint64 {
	var s []byte
	if shift > 0 {
		s = shiftBytesLeft(id[:10], shift)[:8]
	} else {
		s = id[:8]
	}

	return binary.BigEndian.Uint64(s)
}

func shiftBytesLeft(a []byte, l byte) (dst []byte) {
	lb := l / 8
	l = l % 8
	if int(lb) >= len(a) {
		return make([]byte, len(a))
	}
	n := len(a) - int(lb)
	dst = make([]byte, len(a))
	for i := 0; i < n-1; i++ {
		dst[i] = a[i+int(lb)] << l
		dst[i] = (dst[i] & (255 << l)) | (a[i+int(lb)+1] >> (8 - l))
	}
	dst[n-1] = a[n+int(lb)-1] << l
	return dst
}

func setBit(n int, pos uint) int {
	n |= 1 << pos

	return n
}

func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask

	return n
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)

	return val > 0
}
