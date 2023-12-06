package main

import (
	"log"
	"net"
	"net/rpc"
	"net/http"
	"github.com/alfredfo/chord/transport"
)

func serve(tpn *transport.TransportNode) error {
	rpc.Register(tpn)
	rpc.HandleHTTP()

	node := tpn.Node
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
	return nil
}
