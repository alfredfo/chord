package api

import (
	"net"
	"sync"
)

type Key = string
type Value = []byte
type Bucket map[Key]Value
type NodeID = string

const (
	KeySize = 8
)

type NodeInfoType struct {
	ID      NodeID
	TCPAddr net.TCPAddr
}

type Node struct {
	NodeInfo    NodeInfoType
	FingerTable []NodeInfoType
	Predecessor NodeInfoType
	Successors  []NodeInfoType
	Bucket      Bucket
  Backup      Bucket
	Mu          sync.Mutex
}
