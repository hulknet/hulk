package types

import (
	"math"
	"math/bits"
)

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
