package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"github.com/alfredfo/chord/transport"
)

const (
	m = 160
)

var (
	tp                   transport.TransportNode
	stabilizeTime        time.Duration
	fixFingersTime       time.Duration
	checkPredecessorTime time.Duration
	finished             bool
	mutex                sync.Mutex
)

// Don't support custom node id now 
// node id is always hashed from tcp address
func NewNode(addr *net.TCPAddr) (*api.Node, error) {
  
  nodeInfo := api.NodeInfoType{}
	nodeInfo.ID = hashing.HashTcpAddressToString(addr)  
  nodeInfo.TCPAddr = *addr 

	return &api.Node{
	  NodeInfo: nodeInfo,	
		Successor:   api.NodeInfoType{},
		Predecessor: api.NodeInfoType{},
		FingerTable: make([]api.NodeInfoType, m),
		Bucket:      api.Bucket{},
	}, nil

}

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
	tp = transport.TransportNode{}
	var err error
	bindTcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		log.Printf("Failed to resolve tcp address to bind, ip:%v, port:%v, err: %v", bindAddr, bindPort, err)
		return
	}
	log.Printf("manualId: %v", manualID)
	node, err := NewNode(bindTcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	tp.Node = node

	// Output Chord node information
	log.Printf("Chord node ID: %s\n", node.NodeInfo.ID)
	log.Printf("Bind address: %s\n", bindAddr)
	log.Printf("Bind port: %d\n", bindPort)
	go serve(&tp)

	if joinAddr != "" {
		joinTCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", joinAddr, joinPort))
		if err != nil {
			log.Printf("Failed to resolve tcp address to join, ip:%v, port:%v, err: %v", joinAddr, joinPort, err)
			return
		}
    
    // ask the join node to find the sucessor of node args.NodeInfo
  // argsf := FindSuccessorRPCArgs{}
  // replyf := FindSuccessorRPCReply{}

  // err := tp.FindSuccessor(&argsf, &replyf)
  // if err != nil {
  //   log.Println(err)
  // }
    
  // log.Println("finding sucessor first before join...")
    // succ, err := transport.SendFindSuccessor(node, joinTCPAddr)
    // if err != nil {
    //   log.Println(err)
    // }
    // log.Printf("[main] successor found: %v\n", succ)
     

		log.Printf("joining ring\n")
    succ, pred, err := transport.SendJoin(node, joinTCPAddr)
    if err != nil {
      log.Println("transport.SendJoin err: ", err)
    }
    // set successor and predecessor for the current node,
    // since SendJoin only change info at the sucessor node side
    node.Successor = succ
    node.Predecessor = pred

    // tell predecessor to change sucessor to the current node
    log.Printf("Tell node: %v to change sucessor to me: %v, since it's my predecessor now", succ, node.NodeInfo)
    err = transport.SendChangeSucessor(node, &succ.TCPAddr)
    if err != nil {
      log.Println(err)
    }
	} else {
		log.Println("Creating a new ring")
		// if it's a new ring, pionter the predecessor and sucessor to itself
	  node.Successor = node.NodeInfo 
    node.Predecessor = node.NodeInfo
	}

	go stabilizeTimer(stabilizeTime)
	go checkPredecessorTimer(checkPredecessorTime)
	go fixFingersTimer(fixFingersTime)

	cli(bindAddr, bindPort)

	for !finished {
		time.Sleep(time.Second)
	}
}

func MPrintf(format string, args ...interface{}) {
	message := "[main] " + fmt.Sprintf(format, args...)
	log.Print(message)
}
