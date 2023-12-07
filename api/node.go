package api

import (
	"net"
	"sync"
	"math/big"
)

type Key string
type Value string
type Bucket map[Key]Value
type NodeAddress = big.Int
type Node struct {
	ID          NodeAddress
	FingerTable []NodeAddress
	Predecessor NodeAddress
	Successors  []NodeAddress
	Bucket      Bucket
	Address     *net.TCPAddr
	Mu          sync.Mutex
}
