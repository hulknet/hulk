package main

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
	"math/rand"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
)

var (
	nodeCount      int
	nodeCountBLen  int
	maxHop         int
	targetCpl      int
	targetCplDiv   float64
	nodeBatchMulti int
	replications   int
	minReplicas    int
	iteration      int
	nl             []crypto.ID
	nm             map[crypto.ID][]crypto.ID
)

func ParseFlags() error {
	flag.IntVar(&nodeCount, "nodes", 1000, "Nodes in network")
	flag.IntVar(&maxHop, "max-hop", 1, "Maximum network hop")
	flag.IntVar(&targetCpl, "target-cpl", 0, "Target cpl ")
	flag.Float64Var(&targetCplDiv, "target-cpl-div", 2, "Target cpl divider")
	flag.IntVar(&nodeBatchMulti, "node-batch-multi", 2, "Node batch multiplier")
	flag.IntVar(&iteration, "iteration", 9999, "Iterations")
	flag.IntVar(&replications, "rep", 2, "Replications number for transaction")
	flag.IntVar(&minReplicas, "min-rep", 1, "Minimum number of replicas need to find their node")

	flag.Parse()
	return nil
}

func main() {
	err := ParseFlags()
	if err != nil {
		panic(err)
	}

	nodeCountBLen = bits.Len(uint(nodeCount))
	if targetCpl == 0 {
		targetCpl = int(float64(nodeCountBLen) / (math.Sqrt(float64(nodeCountBLen)) / targetCplDiv))
	}
	nl = make([]crypto.ID, nodeCount)
	nm = make(map[crypto.ID][]crypto.ID, nodeCount)

	//create nodes
	for i := 0; i < nodeCount; i++ {
		nl[i] = utils.GenerateSHA()
		nm[nl[i]] = make([]crypto.ID, nodeCountBLen*nodeBatchMulti)
	}

	//foreach node add random connections (nodeCountBLen*2)
	for h := range nm {
		for i := 0; i < nodeCountBLen*nodeBatchMulti; i++ {
			nm[h] = append(nm[h], nl[rand.Intn(nodeCount)])
		}
	}

	fmt.Printf("nodeCount %d \n", nodeCount)
	fmt.Printf("nodeCountBLen %d \n", nodeCountBLen)
	fmt.Printf("node batch size %d \n", nodeCountBLen*nodeBatchMulti)
	fmt.Printf("targetCpl %d \n", targetCpl)
	fmt.Printf("max hop %d \n", maxHop)
	fmt.Printf("Replications %d \n", replications)
	fmt.Printf("Min replicas to match %d \n", minReplicas)
	failure := 0
	for i := 0; i < iteration; i++ {
		if !findClosestPeers() {
			failure++
		}
	}
	fmt.Printf("iterations: %d \n", iteration)
	//fmt.Printf("failure to find: %d \n", failure)
	fmt.Printf("fail rate : %.4f \n", float64(failure)/float64(iteration))
}

func findClosestPeers() bool {
	caller := nl[rand.Intn(nodeCount)]
	target := utils.GenerateSHA()
	repMatch := 0
	for i := 0; i < replications; i++ {
		rep := target.Replica(i)
		if t(caller, rep, maxHop) {
			repMatch++
			if repMatch >= minReplicas {
				return true
			}
		}
	}
	return false
}

func t(id, target crypto.ID, hop int) bool {
	if hop == 0 {
		return false
	}
	for _, nc := range nm[id] {
		cpl := utils.Cpl(target[:], nc[:])
		if cpl > targetCpl || t(nc, target, hop-1) {
			return true
		}
	}
	return false
}
