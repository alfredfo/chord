package main

import (
	"bufio"
	"fmt"
	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/transport"
	"log"
	"net"
	"os"
	"strings"
)

func cli(bindAddr string, bindPort int) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: \n")
		scanner.Scan()
		command := scanner.Text()
		if len(command) != 0 {
			fmt.Println("You entered: ", command)
			splits := strings.Split(command, " ")

			var addr *net.TCPAddr
			var err error
			switch splits[0] {
			case "set":
				key := splits[1]
				val := splits[2]
				if len(splits) > 3 {
					laddr := splits[3]
					lport := splits[4]
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				} else {
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", bindAddr, bindPort))
				}
				if err != nil {
					log.Println(err)
				}
				transport.SendSet(api.Key(key), api.Value(val), addr)
			case "get":
				key := splits[1]
				if len(splits) > 2 {
					laddr := splits[2]
					lport := splits[3]
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				} else {
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", bindAddr, bindPort))
				}
				val, err := transport.SendGet(api.Key(key), addr)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Value for key %v is: %v", key, val)
			case "delete":
				key := splits[1]
				if len(splits) > 2 {
					laddr := splits[2]
					lport := splits[3]
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				} else {
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", bindAddr, bindPort))
				}
				val, err := transport.SendDelete(api.Key(key), addr)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Value for key %v is: %v, deleted from bucket", key, val)
			case "dump": // dump information of a given node
				log.Println("Trying to dump information.")
				if len(splits) > 1 {
					laddr := splits[1]
					lport := splits[2]
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", laddr, lport))
				} else {
					addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", bindAddr, bindPort))
				}
				dumpstr, err := transport.SendDump(addr)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Node info: %v\n", dumpstr)
			default:
				log.Println("not implemented")
			}
		}
	}
}
