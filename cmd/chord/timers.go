package main

import (
	"time"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"github.com/alfredfo/chord/transport"
)

func stabilizeTimer(node *api.Node, ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		node.Mu.Lock()
		defer node.Mu.Unlock()

		x := node.Successor // WRONG ask for pred ... send()

		if hashing.SBetween(node.NodeInfo.ID, x.ID, node.Successor.ID, false) {
			node.Successor = x
		}
		transport.SendNotify(node, node.Successor)
	}
}

func checkPredecessorTimer(node *api.Node, ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		node.Mu.Lock()
		defer node.Mu.Unlock()
	}
}

func fixFingersTimer(node *api.Node, ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
		node.Mu.Lock()
		defer node.Mu.Unlock()
	}
}
