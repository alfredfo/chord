package main

import (
	"log"
	"time"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"github.com/alfredfo/chord/transport"
	"math/big"
)

func stabilizeTimer(node *api.Node, ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		x, err := transport.SendAskPredecessor(&node.Successors[0].TCPAddr)
		// if we can't contact successor then drop it from the list.
		if err != nil || x.ID == "" {
			if len(node.Successors) > 1 {
				node.Successors = node.Successors[1:]
			} else {
				node.Successors = make([]api.NodeInfoType, 0)
				node.Successors = append(node.Successors, node.NodeInfo)
			}
		} else {
			if hashing.SBetween(node.NodeInfo.ID, x.ID, node.Successors[0].ID, true) {
				newSuccs := make([]api.NodeInfoType, 0)
				succs, err := transport.SendAskSuccessors(&x.TCPAddr)
				if err != nil {
					newSuccs = append(newSuccs, node.Successors...)
				} else {
					newSuccs = append(newSuccs, x)
					newSuccs = append(newSuccs, succs...)
				}
				if len(newSuccs) > 4 {
					newSuccs = newSuccs[:3]
				}
				node.Successors = newSuccs
			}
		}
		transport.SendNotify(node, node.Successors[0])
	}
}

func checkPredecessorTimer(node *api.Node, ms int) {
	for !finished {
		// log.Println("===========Check Predecessor==========")
		time.Sleep(time.Millisecond * time.Duration(ms))
		if node.Predecessor.ID == "" {
			continue
		}

		err := transport.SendCheckPredecessor(&node.Predecessor.TCPAddr)

		if err != nil {
			log.Printf("Check predecessor failed: %v, set predecessor to nil.", err)
			node.Predecessor = api.NodeInfoType{}
		}

	}
}

func fixFingersTimer(node *api.Node, ms int, next *int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		*next = *next + 1
		if *next > m {
			*next = 1
		}
		j := hashing.Jump(node.NodeInfo.ID, *next)
		ni, err := transport.SendFindSuccessor(j, &node.Successors[0].TCPAddr)
		if err != nil {
			node.FingerTable[big.NewInt(int64(*next)).String()] = ni
			log.Printf("finger %v\n", node.FingerTable)
		}
	}
}
