package main

import (
	"log"
	"net"
	"net/rpc"
	"net/http"
	"github.com/alfredfo/chord/api"
)


func serve(node *api.Node) error {
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
	return nil
}
