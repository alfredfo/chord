package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)


type Key string
type NodeAddress string

type JoinRPCArgs struct {
	ID NodeAddress
}

type JoinRPCReply struct {
	Ok bool
}

var (
	node            *Node
	stabilizeTime        time.Duration
	fixFingersTime       time.Duration
	checkPredecessorTime time.Duration
	finished             bool
	mutex                sync.Mutex
)

type Node struct {
	ID          NodeAddress
	FingerTable []NodeAddress
	Predecessor NodeAddress
	Successors  []NodeAddress

	Bucket  map[Key]string
	Address *net.TCPAddr
}

func (node *Node) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
  log.Printf("node with ID: %v is joining the ring through: \n", args.ID, node.ID)
	reply.Ok = true
	return nil
}

func (node *Node) findSuccessor(*int) *Node {
	return nil
}

func (node *Node) notify(childNode *Node) {

}

func (node *Node) serve() {
	rpc.Register(node)
	rpc.HandleHTTP()

	addr := node.Address
	log.Println("listening on: ", addr)
	l, e := net.ListenTCP("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func hash(data []byte) string {
	sha1 := sha1.Sum(data)
	s := fmt.Sprintf("%x", sha1)
	return s
}

func hashAddress(tcpAddr net.TCPAddr) NodeAddress {
	return NodeAddress(hash([]byte(tcpAddr.String())))
}

// NewNode creates a new Chord node with the given ID.
func NewNode(id NodeAddress, addr *net.TCPAddr) (*Node, error) {
	if id == "" {
		id = hashAddress(*addr)
	}

	return &Node{
		ID:          id,
		Successors:  nil,
		Predecessor: "",
		FingerTable: make([]NodeAddress, m),
		Bucket:      make(map[Key]string),
		Address:     addr,
	}, nil
}

func call(rpcname string, addr *net.TCPAddr, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", addr.String())
	if err != nil {
		log.Printf("error dialing: %v", err)
		return err
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err != nil {
		log.Printf("error calling: %v", err)
		return err
	}
	return nil
}

const m = 160

func main() {
	// Parse command-line arguments
	var (
		bindAddr             string
		bindPort             int
		joinAddr             string
		joinPort             int
		stabilizeTime        int
		fixFingersTime       int
		checkPredecessorTime int
		numSuccessors        int
		manualID             string
	)

	flag.StringVar(&bindAddr, "a", "0.0.0.0", "The IP address that the Chord client will bind to and advertise.")
	flag.IntVar(&bindPort, "p", 1234, "The port that the Chord client will bind to and listen on.")
	flag.StringVar(&joinAddr, "ja", "", "The IP address of the machine running a Chord node to join.")
	flag.IntVar(&joinPort, "jp", 0, "The port that an existing Chord node is bound to and listening on.")
	flag.IntVar(&stabilizeTime, "ts", 500, "Time in milliseconds between invocations of ‘stabilize’.")
	flag.IntVar(&fixFingersTime, "tff", 500, "Time in milliseconds between invocations of ‘fix fingers’.")
	flag.IntVar(&checkPredecessorTime, "tcp", 500, "Time in milliseconds between invocations of ‘check predecessor’.")
	flag.IntVar(&numSuccessors, "r", 4, "Number of successors maintained by the Chord client.")
	flag.StringVar(&manualID, "i", "", "The identifier (ID) assigned to the Chord client.")

	flag.Parse()

	var err error
	bindTcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
    log.Printf("Failed to resolve tcp address to bind, ip:%v, port:%v, err: %v", bindAddr, bindPort, err)
		return
	}

	var ID NodeAddress = (NodeAddress)(manualID)
	node, err = NewNode(ID, bindTcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	// Output Chord node information
	log.Printf("Chord node ID: %s\n", node.ID)
	log.Printf("Bind address: %s\n", bindAddr)
	log.Printf("Bind port: %d\n", bindPort)
  log.Println("Creating a new ring")
  go node.serve()

	// If joining an existing ring, attempt to join
	if joinAddr != "" {
		joinTcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", joinAddr, joinPort))
		if err != nil {
      log.Printf("Failed to resolve tcp address to join, ip:%v, port:%v, err: %v", joinAddr, joinPort, err)
			return
		}

		args := JoinRPCArgs{}
		reply := JoinRPCReply{}
		args.ID = "3"
		call("Node.Join", joinTcpAddr, &args, &reply)
	}

	for !finished {
		time.Sleep(time.Second)
	}
}
