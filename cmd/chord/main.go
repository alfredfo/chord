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
	"io"
	"net/http"
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
	log.Printf("%v\n", nodeInfo)
	nodeInfo.TCPAddr = *addr

	return &api.Node{
		NodeInfo:    nodeInfo,
		Successors:  make([]api.NodeInfoType, 1),
		Predecessor: api.NodeInfoType{},
		FingerTable: make([]api.NodeInfoType, api.KeySize),
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
	req, err := http.NewRequest("GET", "https://ip.me/", nil)
	req.Header.Add("User-Agent", "curl/8.4.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Printf("could not get public ip via ip.me: %v", err)
	}
	ipb, err := io.ReadAll(resp.Body)
	ip := string(ipb)
	log.Printf("%s", ip)

	tp = transport.TransportNode{}
	bindAddr = ip
	bindTcpAddr := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: bindPort,
	}
	// bindTcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))

	if err != nil {
		log.Printf("Failed to resolve tcp address to bind, ip:%v, port:%v, err: %v", bindAddr, bindPort, err)
		return
	}
	log.Printf("manualId: %v", manualID)
	node, err := NewNode(&bindTcpAddr)

	if err != nil {
		log.Println(err)
		return
	}
	tp.Node = node

	a := hashing.NodeIDToBigInt(node.NodeInfo.ID)
	b := a.String()
	c := hashing.NodeIDToBigInt(b)
	log.Printf("%v %v %v %v\n", node.NodeInfo.ID, a, b, c)

	// Output Chord node information
	log.Printf("Chord node ID: %s\n", node.NodeInfo.ID)
	log.Printf("Bind address: %s\n", bindAddr)
	log.Printf("Bind port: %d\n", bindPort)

	go serve(&tp, bindPort)

	if joinAddr != "" {
		join(node, joinAddr, joinPort)
	} else {
		log.Println("Creating a new ring")
		// if it's a new ring, pionter the predecessor and sucessor to itself
		node.Predecessor = node.NodeInfo
		node.Successors = make([]api.NodeInfoType, 1)
		node.Predecessor = node.Successors[0]
		node.Successors[0] = node.NodeInfo
	}

	go stabilizeTimer(node, stabilizeTime)
	go checkPredecessorTimer(node, checkPredecessorTime)
	var next int
	go fixFingersTimer(node, fixFingersTime, &next)

	cli(bindAddr, bindPort)

	for !finished {
		time.Sleep(time.Second)
	}
}

func join(node *api.Node, joinAddr string, joinPort int) {
	joinTCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", joinAddr, joinPort))
	if err != nil {
		log.Printf("Failed to resolve tcp address to join, ip:%v, port:%v, err: %v", joinAddr, joinPort, err)
		return
	}

	log.Printf("joining ring\n")
	succ, err := transport.SendJoin(node, joinTCPAddr)
	if err != nil {
		log.Println("transport.SendJoin err: ", err)
	}

	kvmap, err := transport.SendCopyData(node.NodeInfo, &succ.TCPAddr)
	log.Println("Received kvmap to store...", kvmap)

	err = transport.SendStoreData(kvmap, &node.NodeInfo.TCPAddr)

	if err != nil {
		log.Println(err)
	}
	// TODO when to delete the key value pair from successor after we stored

	log.Printf("asd")

	node.Successors[0] = succ
	succs, err := transport.SendAskSuccessors(&node.Successors[0].TCPAddr)
	node.Successors = append(node.Successors, succs...)
	if len(node.Successors) > 3 {
		node.Successors = node.Successors[:3]
	}

	//transport.SendChangeSucessor(succ, &succ.TCPAddr)
	log.Printf("asd2")
	// set successor and predecessor for the current node,
	// since SendJoin only change info at the sucessor node side
	node.Successors[0] = succ
	node.Predecessor = node.NodeInfo

	// delete old data from successor
	err2 := transport.SendDeleteOldData(node.NodeInfo, &succ.TCPAddr)
	if err2 != nil {
		log.Println(err)
	}

}

func MPrintf(format string, args ...interface{}) {
	message := "[main] " + fmt.Sprintf(format, args...)
	log.Print(message)
}
