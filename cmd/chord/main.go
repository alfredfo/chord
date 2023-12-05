package main

import (
	"flag"
	"net"
	"time"
	"sync"
	"fmt"
	"log"
	"crypto/sha1"
)

// ChordNode represents a node in the Chord DHT.
type ChordNode struct {
	ID                   string
	Successor            *ChordNode
	Predecessor          *ChordNode
	FingerTable          []*ChordNode
	Data                 map[string]string
	Address              net.TCPAddr	
}

var (
	localNode			*ChordNode
	stabilizeTime        time.Duration
	fixFingersTime       time.Duration
	checkPredecessorTime time.Duration
	mutex                sync.Mutex
)

func (node *ChordNode) findSuccessor(*int) *ChordNode {
	return nil
}

func (node *ChordNode) notify(childNode *ChordNode) {

}

func (node *ChordNode) join() *ChordNode {
	return nil
}


func hash(tcpAddr net.TCPAddr) string {
	sha1 := sha1.Sum([]byte(tcpAddr.String()))
    s := fmt.Sprintf("%x", sha1)
	return s
}

// NewChordNode creates a new Chord node with the given ID.
func NewChordNode(id string, addr string, port int ) (*ChordNode, error) {

	if id == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}
	return &ChordNode{
		ID:                   id,
		Successor:            nil,
		Predecessor:          nil,
		FingerTable:          make([]*ChordNode, m),
		Data:                 make(map[string]string),
		Address: 			  *tcpAddr,
		
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
		ID             		 string
	)

	flag.StringVar(&bindAddr, "a", "0.0.0.0", "The IP address that the Chord client will bind to and advertise.")
	flag.IntVar(&bindPort, "p", 1234, "The port that the Chord client will bind to and listen on.")
	flag.StringVar(&joinAddr, "ja", "", "The IP address of the machine running a Chord node to join.")
	flag.IntVar(&joinPort, "jp", 0, "The port that an existing Chord node is bound to and listening on.")
	flag.IntVar(&stabilizeTime, "ts", 500, "Time in milliseconds between invocations of ‘stabilize’.")
	flag.IntVar(&fixFingersTime, "tff", 500, "Time in milliseconds between invocations of ‘fix fingers’.")
	flag.IntVar(&checkPredecessorTime, "tcp", 500, "Time in milliseconds between invocations of ‘check predecessor’.")
	flag.IntVar(&numSuccessors, "r", 4, "Number of successors maintained by the Chord client.")
	flag.StringVar(&ID, "i", "", "The identifier (ID) assigned to the Chord client.")

	flag.Parse()



	// Create Chord node
	// var node *ChordNode
	if ID == "" {
		ID = "abc"
	} 
	var err error
	localNode, err = NewChordNode(ID, bindAddr, bindPort)
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
	if createRing{
	} else {
		joinNode, err := NewChordNode(nil, joinAddr , joinPort)
		if err != nil {	
			log.Println(err)
			return
		}
	}

	// Keep the main goroutine running
}
