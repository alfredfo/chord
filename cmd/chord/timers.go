package main

import (
	"log"
	"time"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"github.com/alfredfo/chord/transport"
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
				log.Println("======== set succ to self ======== ")
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
				if len(newSuccs) > 3 {
					newSuccs = newSuccs[:3]
				}
				node.Successors = newSuccs
			}
		}
		transport.SendNotify(node, node.Successors[0])
	}
}

func checkPredecessorTimer(node *api.Node, ms int) {
  var predecessorFixed = false
	for !finished {
		// log.Println("===========Check Predecessor==========")
		time.Sleep(time.Millisecond * time.Duration(ms))
		if node.Predecessor.ID == "" {
			continue
		}

		err := transport.SendCheckPredecessor(&node.Predecessor.TCPAddr)

		if err != nil {
			log.Printf("Check predecessor failed: %v, set predecessor to nil and backup data.", err)
			node.Predecessor = api.NodeInfoType{}
       
      for k, v := range node.Backup {
        node.Bucket[k] = v
      }
      transport.SendBackup(node.Bucket, &node.Successors[0].TCPAddr)
      predecessorFixed = true
		} else {
      if predecessorFixed {
        predecessorFixed = false 
        // TODO backup predecessor
        predBucket, err := transport.SendAskBucket(&node.Predecessor.TCPAddr) 
        if err == nil {
          node.Backup = predBucket
        } else {
          log.Println("err ")
          node.Backup = api.Bucket{}
        }
      }
    }
	}
}

func fixFingersTimer(node *api.Node, ms int, next *int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		for n := 1; n <= api.KeySize; n++ {
			j := hashing.Jump(node.NodeInfo.ID, n)
			ni, err := transport.SendFindSuccessor(j, &node.Successors[0].TCPAddr)
			if err == nil {
				node.FingerTable[n-1] = ni
			}
		}
	}
}
