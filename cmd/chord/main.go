package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
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

var (
	localNode            *Node
	stabilizeTime        time.Duration
	fixFingersTime       time.Duration
	checkPredecessorTime time.Duration
	mutex                sync.Mutex
)

func (node *Node) findSuccessor(*int) *Node {
	return nil
}

func (node *Node) notify(childNode *Node) {

}

func (node *Node) join() *Node {
	return nil
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
		id = hashAddress(*tcpAddr)
	}

	return &Node{
		ID:          id,
		Successors:  nil,
		Predecessor: "",
		FingerTable: make([]NodeAddress, m),
		Bucket:      make(map[Key]string),
		Address:     tcpAddr,
	}, nil
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		return
	}
	var ID NodeAddress = (NodeAddress)(manualID)
	localNode, err = NewNode(ID, tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	// Output Chord node information
	fmt.Printf("Chord Node ID: %s\n", localNode.ID)
	fmt.Printf("Bind Address: %s\n", bindAddr)
	fmt.Printf("Bind Port: %d\n", bindPort)

	createRing := joinAddr != ""
	// If joining an existing ring, attempt to join
	if createRing {
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
		if err != nil {
			return
		}
		joinNode, err := NewNode("", tcpAddr)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Keep the main goroutine running
}
