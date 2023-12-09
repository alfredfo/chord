package main

import (
	"log"
	"time"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"github.com/alfredfo/chord/transport"
)

func stabilizeTimer(node *api.Node, ms int) {
	// log.Printf("fin")
	for !finished {
		// log.Printf("fina")
		time.Sleep(time.Millisecond * time.Duration(ms))

		emptyPred := api.NodeInfoType{}
		succs, pred, err := transport.SendAskPredecessor(&node.Successor.TCPAddr)
		predIsEmpty := pred.ID == emptyPred.ID

		if err != nil {
			node.Successor = node.NodeInfo
			log.Printf("error asking for predecessor in stabilize %v", err)
		} else {
      if predIsEmpty {
        continue
      }
			x := pred
			//newSuccs := make([]api.NodeInfoType, 4)
			//newSuccs = append(newSuccs, node.Successor)
			//newSuccs = append(newSuccs, succs[:2]...)
			//node.Successors = newSuccs
			node.Successor = succs[0]

			if hashing.SBetween(node.NodeInfo.ID, x.ID, node.Successor.ID, false) {
				node.Successor = x
			}
			// log.Printf("stab %v", pred)
			transport.SendNotify(node, node.Successor)
		}
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

func fixFingersTimer(node *api.Node, ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		// node.Mu.Lock()
		// defer node.Mu.Unlock()
	}
}
