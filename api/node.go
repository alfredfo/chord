package api

import (
	"math/big"
	"net"
	"sync"
)

type Key = string
type Value = string
type Bucket map[Key]Value
type NodeAddress = big.Int
type Node struct {
	ID          NodeAddress
	FingerTable []NodeAddress
	Predecessor *Node
	Successor   *Node
	Bucket      Bucket
	Address     *net.TCPAddr
	Mu          sync.Mutex
}
