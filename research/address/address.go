package main

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
	"math/rand"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
	"github.com/montanaflynn/stats"
)

var (
	nodeCount         int
	bucketSize        int
	shaRepNum         int
	nodeBucketMaxSize int
	nl                []crypto.ID
	nm                map[string]int
)

func ParseFlags() {
	flag.IntVar(&nodeCount, "nodes", 2221000, "Nodes in network")
	flag.IntVar(&bucketSize, "bucket-size", 4, "bucket size of bits")
	flag.IntVar(&shaRepNum, "sha-replicas", 0, "number of sha256 of key")

	flag.Parse()
}

func main() {
	ParseFlags()

	nl = make([]crypto.ID, nodeCount)
	nm = make(map[string]int)
	keyLen := bits.Len(uint(nodeCount)) - bucketSize
	for i := 0; i < nodeCount; i++ {
		nl[i] = utils.GenerateSHA()
		if shaRepNum > 0 {
			nl[i] = nl[i].Replica(rand.Intn(shaRepNum))
		}
		key := nl[i].HexL(keyLen)
		if _, ok := nm[key]; ok {
			nm[key]++
			if nodeBucketMaxSize < nm[key] {
				nodeBucketMaxSize = nm[key]
			}
		} else {
			nm[key] = 1
		}
	}
	rawData := make([]int, 0)
	for _, v := range nm {
		rawData = append(rawData, v)
	}

	data := stats.LoadRawData(rawData)
	min, _ := stats.Min(data)
	median, _ := stats.Median(data)
	perc90, _ := stats.Percentile(data, 90.0)
	perc95, _ := stats.Percentile(data, 95.0)
	perc99, _ := stats.Percentile(data, 99.0)
	perc100, _ := stats.Percentile(data, 100.0)
	medianDeviation, _ := stats.MedianAbsoluteDeviation(data)
	standardDeviation, _ := stats.StandardDeviation(data)

	fmt.Printf("nodeCount: %d, ", nodeCount)
	fmt.Printf("bucketSize bits: %d, ", bucketSize)
	fmt.Printf("keyLen bits: %d, ", keyLen)
	fmt.Printf("nodeCount bits length: %d, \n", bits.Len(uint(nodeCount)))

	fmt.Printf("nodebuckets count: %d, ", len(nm))
	fmt.Printf("target: %d \n", int(math.Pow(2, float64(keyLen))))

	fmt.Printf("nodebucket size deviation median: %.1f, ", medianDeviation)
	fmt.Printf("standard: %.1f \n", standardDeviation)

	fmt.Printf("nodebucket size min: %.1f,  median: %.1f, p90: %.1f, p95: %.1f, p99: %.1f, p100: %.1f \n",
		min, median, perc90, perc95, perc99, perc100)
}
