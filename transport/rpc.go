package transport

import (
	"log"
	"net"
	"net/rpc"
)

func call(rpcName string, addr *net.TCPAddr, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", addr.String())
	if err != nil {
		log.Printf("error dialing: %v", err)
		return err
	}
	defer c.Close()

	err = c.Call(rpcName, args, reply)
	if err != nil {
		log.Printf("error calling: %v", err)
		return err
	}
	return nil
}
