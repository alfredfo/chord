package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/alfredfo/chord/api"
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

const keySize = 4

func hash(data []byte) string {
	sha1 := sha1.Sum(data)
	s := fmt.Sprintf("%x", sha1)
	return s
}

func hashStringToBigInt(elt string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(elt))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

func hashAddress(tcpAddr net.TCPAddr) api.NodeAddress {
	hashInt := hashStringToBigInt(tcpAddr.String())
	// hashMod := new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)
	nodeId := new(big.Int).Mod(hashInt, big.NewInt(keySize))
	log.Printf("nodeid: %v", nodeId)
	return api.NodeAddress(nodeId.String())
}

// Newapi.Node creates a new Chord node with the given ID.
func NewNode(id api.NodeAddress, addr *net.TCPAddr) (*api.Node, error) {
	if id == "" {
		id = hashAddress(*addr)
	}

	return &api.Node{
		ID:          id,
		Successors:  nil,
		Predecessor: "",
		FingerTable: make([]api.NodeAddress, m),
		Bucket:      api.Bucket{},
		Address:     addr,
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

	var ID api.NodeAddress = (api.NodeAddress)(manualID)
	node, err := NewNode(ID, bindTcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	ID = node.ID
	tp.Node = node
	// Output Chord node information
	log.Printf("Chord node ID: %s\n", node.ID)
	log.Printf("Bind address: %s\n", bindAddr)
	log.Printf("Bind port: %d\n", bindPort)
	go serve(&tp)

	if joinAddr != "" {
		joinTCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", joinAddr, joinPort))
		if err != nil {
			log.Printf("Failed to resolve tcp address to join, ip:%v, port:%v, err: %v", joinAddr, joinPort, err)
			return
		}
		log.Printf("joining ring\n")
		transport.SendJoin(node.ID, joinTCPAddr)

	} else {
		log.Println("Creating a new ring")
	}

	for !finished {
		time.Sleep(time.Second)
	}
}
