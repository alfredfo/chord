package api

import (
	"net"
	"sync"
)

type Key string
type Value string
type Bucket map[Key]Value
type NodeAddress string
type Node struct {
	ID          NodeAddress
	FingerTable []NodeAddress
	Predecessor NodeAddress
	Successors  []NodeAddress
	Bucket      Bucket
	Address     *net.TCPAddr
	Mu          sync.Mutex
}
