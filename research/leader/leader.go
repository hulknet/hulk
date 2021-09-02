package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"

	"github.com/kotfalya/hulk/pkg/utils"
)

type Link struct {
	id       sha
	cpl      int
	distance *big.Int
}

type sha [32]byte

type Node struct {
	id    sha
	links map[sha]Link
	l     []Link
}

func (n Node) AddLink(l Link) {

}

func (n sha) ID() string {
	return hex.EncodeToString(n[:])
}

var (
	nodeNumbers      int
	nodeNumbersLen   int
	biggestCpl       int
	targetBiggestCpl int
	targetCplCount   int
	targetModCount   int
	shortestDist     *big.Int
)

func main() {
	nodeNumbers = utils.Random()
	if nodeNumbers < 100 {
		nodeNumbers = 100
	}
	if nodeNumbers > 1000000 {
		nodeNumbers = 1000000
	}
	nodeNumbersLen = bits.Len(uint(nodeNumbers))
	nL := big.NewInt(int64(nodeNumbers))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
	targetHash := sha256.Sum256(bs)

	fmt.Println("nodes")
	fmt.Println(nodeNumbers)

	nl := make([]sha, nodeNumbers)
	nm := make(map[sha]Node, nodeNumbers)
	for i := 0; i < nodeNumbers; i++ {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(utils.Random()*i))
		nodeHash := sha256.Sum256(bs)
		nm[nodeHash] = Node{nodeHash, make(map[sha]Link, 0), make([]Link, 0)}
		nl[i] = nodeHash

		cpl := utils.Cpl(nodeHash[:], targetHash[:])
		if cpl >= nodeNumbersLen/2 {
			targetCplCount++
		}
		nD := big.NewInt(0).SetBytes(nodeHash[:])

		if big.NewInt(0).Mod(nD, nL).String() == "0" {
			//fmt.Println(big.NewInt(0).Mod(nD, nL).String())
			//fmt.Println(nD.String())
			//fmt.Println(nL.String())
			targetModCount++
		}
		if cpl > targetBiggestCpl {
			targetBiggestCpl = cpl
		}
	}

	for _, n := range nm {
		for {
			if len(n.links) == 20 {
				break
			}

			l := nl[rand.Intn(nodeNumbers)]
			if _, ok := n.links[l]; ok || l == n.id {
				continue
			}
			cpl := utils.Cpl(n.id[:], l[:])
			distance := utils.Distance(n.id[:], l[:])
			if cpl > biggestCpl {
				biggestCpl = cpl
				shortestDist = distance
			}
			n.links[l] = Link{
				l,
				cpl,
				distance,
			}
		}
	}

	//for _, l := range nm[nl[99]].links {
	//	fmt.Println(l.id.Hex())
	//	fmt.Println(l.cpl)
	//	fmt.Println(l.distance)
	//	fmt.Println("----------")
	//}

	fmt.Println("nodes")
	fmt.Println(nodeNumbers)
	fmt.Println("nodes len")
	fmt.Println(bits.Len(uint(nodeNumbers)))
	fmt.Println("biggest cpl")
	fmt.Println(biggestCpl)
	fmt.Println("target cpl count")
	fmt.Println(targetCplCount)
	fmt.Println("target biggest cpl")
	fmt.Println(targetBiggestCpl)
	fmt.Println("target mod count")
	fmt.Println(targetModCount)
	fmt.Println("shortest distance")
	fmt.Println(shortestDist)
}
