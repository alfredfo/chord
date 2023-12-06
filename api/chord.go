package api

import (
	"net"
	"log"
	"net/http"
	"net/rpc"
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

func (node *Node) FindSuccessor(*int) *Node {
	return nil
}

func (node *Node) Notify(childNode *Node) {

}

func (node *Node) Serve() {
	rpc.Register(node)
	rpc.HandleHTTP()

	addr := node.Address
	log.Println("listening on: ", addr)
	l, e := net.ListenTCP("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	err := http.Serve(l, nil)
	if err != nil {
		log.Fatal(err)
	}
}
