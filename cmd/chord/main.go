package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

// Newapi.Node creates a new Chord node with the given ID.
func NewNode(sid string, addr *net.TCPAddr) (*api.Node, error) {
	var id *api.NodeAddress
	if sid == "" {
		id = hashAddress(addr)
	}

	node := &api.Node{
		ID:          *id,
		Successor:   nil,
		Predecessor: nil,
		FingerTable: make([]api.NodeAddress, m),
		Bucket:      api.Bucket{},
		Address:     addr,
	}
	
	return node, nil
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

	node, err := NewNode(manualID, bindTcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	tp.Node = node

	// Output Chord node information
	log.Printf("Chord node ID: %s\n", node.ID.String())
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
		node.Successor = node
	}

	go stabilizeTimer(stabilizeTime)
	go checkPredecessorTimer(checkPredecessorTime)
	go fixFingersTimer(fixFingersTime)
	
	cli(bindAddr, bindPort)

	for !finished {
		time.Sleep(time.Second)
	}
}

func cli(bindAddr string, bindPort int) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: \n")
		// reads user input until \n by default
		scanner.Scan()
		// Holds the string that was scanned
		command := scanner.Text()
		if len(command) != 0 {
			fmt.Println("You entered: ", command)
			// Here you can add a switch or if statements to handle the commands
			splits := strings.Split(command, " ")
			switch splits[0] {
			case "set":
				key := splits[1]
				val := splits[2]
				var addr *net.TCPAddr
				var err error
				if len(splits) > 3 {
					laddr := splits[3]
					lport := splits[4]
				
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				} else {
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", "localhost", bindPort))
				}
				if err != nil {
					log.Println(err)
				}
				transport.SendSet(api.Key(key), api.Value(val), addr)
			case "get":
				key := splits[1]
				laddr := splits[2]
				lport := splits[3]
				addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				val, err := transport.SendGet(api.Key(key), addr)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Value for key %v is: %v", key, val)
			default:
				log.Println("not implemented")
			}
		} else {
			// exit if user entered an empty string
			break
		}
	}
	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}
