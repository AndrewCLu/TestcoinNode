package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/node"
)

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/AndrewCLu/TestcoinNode/block"
// 	"github.com/AndrewCLu/TestcoinNode/crypto"
// 	"github.com/AndrewCLu/TestcoinNode/util"
// )

// func testSolveSpeed(targetFirstFour [4]byte, numIters int) {
// 	var target [crypto.HashLength]byte
// 	target[0] = targetFirstFour[0]
// 	target[1] = targetFirstFour[1]
// 	target[2] = targetFirstFour[2]
// 	target[3] = targetFirstFour[3]
// 	for i := 4; i < 32; i++ {
// 		target[i] = byte(255)
// 	}

// 	var totalTime int64 = 0
// 	var maxTime int64 = 0
// 	var minTime int64 = 0
// 	for i := 0; i < numIters; i++ {
// 		var random1 [crypto.HashLength]byte
// 		s1 := rand.NewSource(time.Now().UnixNano())
// 		r1 := rand.New(s1)
// 		fill1 := make([]byte, crypto.HashLength)
// 		r1.Read(fill1)
// 		copy(random1[:], fill1)

// 		var random2 [crypto.HashLength]byte
// 		s2 := rand.NewSource(time.Now().UnixNano())
// 		r2 := rand.New(s2)
// 		fill2 := make([]byte, crypto.HashLength)
// 		r2.Read(fill2)
// 		copy(random2[:], fill2)

// 		header := block.BlockHeader{
// 			ProtocolVersion:     uint16(1),
// 			PreviousBlockHash:   random1,
// 			AllTransactionsHash: random2,
// 			Timestamp:           time.Now().Round(0),
// 			Target:              target,
// 			Nonce:               uint32(0),
// 		}

// 		runTime := header.Solve()
// 		totalTime += runTime

// 		if maxTime == int64(0) || runTime > maxTime {
// 			maxTime = runTime
// 		}

// 		if minTime == int64(0) || runTime < minTime {
// 			minTime = runTime
// 		}
// 	}
// 	fmt.Printf("Solved difficulty %v with average time %v seconds. Maximum solve time was %v seconds, minimum time was %v seconds.\n", util.HashToHexString(target), totalTime/int64(numIters), maxTime, minTime)
// }

func main() {
	fmt.Println("\n\n\n-------------------BEGINNING TEST-------------------")
	node, _ := node.New()

	satoshi := node.NewAccount()
	node.Initialize(satoshi.Address)
	node.BeginMiner(satoshi.Address)

	alice := node.NewAccount()
	bob := node.NewAccount()

	fmt.Println("Creating the blockchain...")
	node.PrintChainState()

	node.NewPeerTransaction(satoshi, alice.Address, 3, 1)
	node.MineBlock()

	node.NewPeerTransaction(satoshi, bob.Address, 4, 2)
	node.NewPeerTransaction(satoshi, bob.Address, 4, 2.4)
	node.NewPeerTransaction(satoshi, bob.Address, 4, 2.3)
	node.NewPeerTransaction(satoshi, bob.Address, 4, 2.2)
	node.NewPeerTransaction(satoshi, bob.Address, 4, 2.1)
	node.NewPeerTransaction(satoshi, bob.Address, 4, 2.05)
	node.MineBlock()

	fmt.Println("Distributing the wealth...")
	node.PrintChainState()

	node.NewPeerTransaction(alice, bob.Address, .5, .25)
	node.MineBlock()

	fmt.Println("Alice owes bob...")
	node.PrintChainState()
}
