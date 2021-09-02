package router

import "math/bits"

const (
	TargetBucketSizeOffset = 2
	BucketMaxSize          = 8
)

type BucketAddr uint8

type Addr uint64

type AddrMap struct {
	addr       Addr
	bucketNum  int
	bucketSize int
}

func newAddrMap(addr Addr, bitSize int) AddrMap {
	bucketNum := ((bitSize + TargetBucketSizeOffset) / BucketMaxSize) + 1
	bucketSize := (bitSize + TargetBucketSizeOffset) / bucketNum
	return AddrMap{addr, bucketNum, bucketSize}
}

func (a AddrMap) BucketAddr(addr Addr) BucketAddr {
	leftShift := a.bucketSize * a.BucketIndex(addr)
	rightShift := 64 - leftShift - a.bucketSize
	return BucketAddr((addr << leftShift) >> rightShift)
}

func (a AddrMap) Cpl(addr Addr) int {
	return bits.LeadingZeros64(uint64(a.addr ^ addr))
}

func (a AddrMap) BucketNum() int {
	return a.bucketNum
}

func (a AddrMap) BucketSize() int {
	return a.bucketSize
}

func (a AddrMap) BucketIndex(addr Addr) int {
	cpl := a.Cpl(addr)
	for i := 0; i < a.bucketNum; i++ {
		if cpl <= a.bucketSize*(i+1) {
			return i
		}
	}
	return a.bucketNum - 1
}
