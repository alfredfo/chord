package api

import (
	"net"
	"sync"
)

type Key = string
type Value = string
type Bucket map[Key]Value
type NodeID = string

type NodeInfoType struct {
	ID      NodeID
	TCPAddr net.TCPAddr
}

type Node struct {
	NodeInfo    NodeInfoType
	FingerTable map[NodeID]NodeInfoType
	Predecessor NodeInfoType
	Successors  []NodeInfoType
	Bucket      Bucket
	Mu          sync.Mutex
}
