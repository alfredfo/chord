package main

import (
	"github.com/alfredfo/chord/transport"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// func serve(tpn *transport.TransportNode) error {

// 	if err := rpc.Register(tpn); err != nil {
// 		log.Fatal("Failed to register RPC method")
// 	}
// 	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
// 	if err != nil {
// 		log.Fatalf("server: loadkeys: %s", err)
// 	}
// 	if len(cert.Certificate) != 2 {
// 		log.Fatal("server.crt should have 2 concatenated certificates: server + CA")
// 	}
// 	ca, err := x509.ParseCertificate(cert.Certificate[1])
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	certPool := x509.NewCertPool()
// 	certPool.AddCert(ca)
// 	config := tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		ClientAuth:   tls.RequireAndVerifyClientCert,
// 		ClientCAs:    certPool,
// 	}
// 	config.Rand = rand.Reader
// 	node := tpn.Node
// 	addr := node.NodeInfo.TCPAddr

// 	service := addr.String()
// 	listener, err := tls.Listen("tcp", service, &config)
// 	if err != nil {
// 		log.Fatalf("server: listen: %s", err)
// 	}

// 	rpc.Register(tpn)

// 	log.Println("listening on: ", addr)
// 	l, e := net.ListenTCP("tcp", &addr)
// 	if e != nil {
// 		log.Fatal("listen error:", e)
// 	}
// 	err := http.ServeTLS(l, nil, "cert.crt", "key.key")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func serve(tpn *transport.TransportNode) error {
	rpc.Register(tpn)
	rpc.HandleHTTP()

	node := tpn.Node
	addr := node.NodeInfo.TCPAddr
	log.Println("listening on: ", addr)
	l, e := net.ListenTCP("tcp", &addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	err := http.ServeTLS(l, nil, "certs/signed.crt", "certs/cert.key")
	if err != nil {
		log.Fatalf("error in ServeTLS: %v\n", err)
	}
	return nil
}
