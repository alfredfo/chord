package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"

	"github.com/alfredfo/chord/transport"
)

func serve(tpn *transport.TransportNode, bindPort int) error {
	if err := rpc.Register(tpn); err != nil {
		log.Fatal("Failed to register RPC method")
	}
	cert, err := tls.LoadX509KeyPair("certs/signed.crt", "certs/cert.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		Rand:         rand.Reader,
	}
	listener, err := tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", bindPort), config)
	if err != nil {
		log.Fatalf("error in tls.Listen: %v\n", err)
	}
	defer listener.Close()
	// log.Println("RPC server is running on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(conn)
	}
}

//func serve(tpn *transport.TransportNode) error {
//	rpc.Register(tpn)
//	rpc.HandleHTTP()
//
//	node := tpn.Node
//	addr := node.NodeInfo.TCPAddr
//	log.Println("listening on: ", addr)
//	l, e := net.ListenTCP("tcp", &addr)
//	if e != nil {
//		log.Fatal("listen error:", e)
//	}
//	err := http.ServeTLS(l, nil, "certs/signed.crt", "certs/cert.key")
//	if err != nil {
//		log.Fatalf("error in ServeTLS: %v\n", err)
//	}
//	return nil
//}
