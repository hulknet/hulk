package types

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
	rd "math/rand"
	"time"
)

func Random() int {
	rd.Seed(time.Now().UnixNano())
	return rd.Int()
}

func GenerateSHA() ID {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(Random()))
	return sha256.Sum256(bs)
}

func GenerateSHAFrom(source string) ID {
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
			return i*8 + bits.LeadingZeros8(b)
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

func Normalize(source uint64, bitSizePrefix uint8) float64 {
	val := clearBitPrefix(source, bitSizePrefix)
	max := clearBitPrefix(math.MaxUint64, bitSizePrefix)
	return float64(val) / float64(max)
}

func clearBitPrefix(val uint64, bitSizePrefix uint8) uint64 {
	return (val << bitSizePrefix) >> bitSizePrefix
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
