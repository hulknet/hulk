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
	targetValue int
	minValue    int
	ledgerID    crypto.ID
)

func ParseFlags() error {
	flag.IntVar(&targetValue, "target-value", 20, "Min ones in xor result")
	flag.IntVar(&minValue, "min-value", 250, "Start value for ones in xor result")
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
	fmt.Printf("targetValue %d \n", targetValue)
	fmt.Printf("start: %s \n", t1)
	var iteration int
	for {
		iteration++
		if success := election(iteration); success {
			fmt.Printf("election succeed time: %s, iteration %d \n", time.Now().Sub(t1).String(), iteration)
			return
		}
	}
}

func election(iteration int) bool {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(utils.Random()))
	txID := sha256.Sum256(bs)

	value := utils.OneCount(txID[:], ledgerID[:])

	if minValue > value {
		minValue = value
		fmt.Printf("minValue is: %d, iteration %d \n", minValue, iteration)
	}

	if value <= targetValue {
		return true
	}

	return false
}
