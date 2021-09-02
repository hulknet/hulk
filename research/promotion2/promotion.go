package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"time"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
)

var (
	targetCpl  int
	biggestCPl int
	ledgerID   crypto.ID
)

func ParseFlags() error {
	flag.IntVar(&targetCpl, "target-cpl", 20, "Minimum common prefix length")

	flag.Parse()
	return nil
}

func main() {
	err := ParseFlags()
	if err != nil {
		panic(err)
	}

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
	ledgerID = sha256.Sum256(bs)

	t1 := time.Now()
	fmt.Printf("targetCpl %d \n", targetCpl)
	fmt.Printf("start: %s \n", t1)
	var iteration int
	for {
		iteration++
		if success := election(iteration); success {
			fmt.Printf("election succeed time: %s, iteration %d \n", time.Now().Sub(t1).String(), iteration)
			t1 = time.Now()
			iteration = 0
		}
	}
}

func election(iteration int) bool {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
	txID := sha256.Sum256(bs)

	cpl := utils.Cpl(txID[:15], ledgerID[:15])

	if biggestCPl < cpl {
		biggestCPl = cpl
		fmt.Printf("biggestCPl is: %d, iteration %d \n", biggestCPl, iteration)
	}

	if cpl >= targetCpl {
		return true
	}

	return false
}
