package api

import (
	"net"
)

type Key string
type NodeAddress string
type Node struct {
	ID          NodeAddress
	FingerTable []NodeAddress
	Predecessor NodeAddress
	Successors  []NodeAddress

	Bucket  map[Key]string
	Address *net.TCPAddr
}
