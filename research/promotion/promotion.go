package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
)

var (
	nodeCount     int
	nodeCountBLen int
	minCpl        int
	targetCpl     int
	nl            []crypto.ID
	biggestCPl    int
)

func ParseFlags() error {
	flag.IntVar(&nodeCount, "nodes", 1000, "Nodes in network")
	flag.IntVar(&minCpl, "min-cpl", 20, "Minimum common prefix length")

	flag.Parse()
	return nil
}

func main() {
	err := ParseFlags()
	if err != nil {
		panic(err)
	}

	nodeCountBLen = bits.Len(uint(nodeCount))
	targetCpl = nodeCountBLen + int(math.Sqrt(float64(nodeCountBLen)))
	if targetCpl < minCpl {
		targetCpl = minCpl
	}
	nl = make([]crypto.ID, nodeCount)
	for i := 0; i < nodeCount; i++ {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
		nodeHash := sha256.Sum256(bs)
		nl[i] = nodeHash
	}

	t1 := time.Now()
	fmt.Printf("nodeCount %d \n", nodeCount)
	fmt.Printf("nodeCountBLen %d \n", nodeCountBLen)
	fmt.Printf("targetCpl %d \n", targetCpl)
	fmt.Printf("start: %s \n", t1)
	var iteration int
	for {
		iteration++
		if success := election(iteration); success {
			fmt.Printf("election success time: %s, iteration %d \n", time.Now().Sub(t1).String(), iteration)
			t1 = time.Now()
			iteration = 0
		}
	}
}

func election(iteration int) bool {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
	txID := sha256.Sum256(bs)
	//ncL := big.NewInt(int64(nodeCount))

	txHolder := make(map[crypto.ID]bool, 0)
	for {
		if len(txHolder) == nodeCountBLen {
			break
		}

		nID := nl[rand.Intn(nodeCount)]
		if _, ok := txHolder[nID]; ok {
			continue
		}
		txHolder[nID] = true
		pID := sha256.Sum256(append(txID[:], nID[:]...))
		cpl := utils.Cpl(pID[:nodeCountBLen], nID[:nodeCountBLen])
		if biggestCPl < cpl {
			biggestCPl = cpl
		}

		if cpl >= targetCpl {
			return true
		}
	}
	return false
}
