package router

type Table struct {
	buckets []*Bucket
	addrMap AddrMap
}

func newTable(addr Addr, bitSize int) *Table {
	addrMap := newAddrMap(addr, bitSize)
	var buckets []*Bucket
	for i := 0; i < addrMap.bucketNum; i++ {
		buckets = append(buckets, newBucket(bitSize))
	}

	return &Table{buckets, addrMap}
}

func (t Table) GetPeer(addr Addr) (Peer, bool) {
	return t.getBucket(addr).Get(t.addrMap.BucketAddr(addr))
}

func (t Table) SetPeer(peer Peer) {
	peer.BucketAddr = t.addrMap.BucketAddr(peer.Addr)
	t.getBucket(peer.Addr).Set(peer)
}

func (t Table) BitMask() []uint64 {
	var mask []uint64
	for _, b := range t.buckets {
		mask = append(mask, b.bitMask)
	}

	return mask
}

func (t Table) getBucket(addr Addr) *Bucket {
	return t.buckets[t.addrMap.BucketIndex(addr)]
}
