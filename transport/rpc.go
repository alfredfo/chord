package transport

import (
	"crypto/tls"
	"log"
	"net"
	"net/rpc"
)

func call(rpcName string, addr *net.TCPAddr, args interface{}, reply interface{}) error {
	cert, err := tls.LoadX509KeyPair("certs/combined.crt", "certs/cert.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true, // For testing purposes only, consider removing this in production
	}
  conn, err := tls.Dial("tcp", addr.String(), config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}

	defer conn.Close()

	// log.Println("client: connected to: ", conn.RemoteAddr())
	rpcClient := rpc.NewClient(conn)

	if err := rpcClient.Call(rpcName, args, reply); err != nil {
		log.Fatal("Failed to call RPC: ", err)
	}
	return nil
}
