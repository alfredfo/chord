package transport

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
)

func call(rpcName string, addr *net.TCPAddr, args interface{}, reply interface{}) error {
	cert, err := tls.LoadX509KeyPair("certs/combined.crt", "certs/cert.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	if len(cert.Certificate) != 2 {
		log.Fatal("client.crt should have 2 concatenated certificates: client + CA")
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	// config := tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// 	RootCAs:      certPool,
	// }
	config := tls.Config{
		RootCAs: certPool,
	}
	conn, err := tls.Dial("tcp", addr.String(), &config)

	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}

	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())
	rpcClient := rpc.NewClient(conn)

	if err := rpcClient.Call(rpcName, args, reply); err != nil {
		log.Fatal("Failed to call RPC: ", err)
	}
	return nil
}
