package api

import (
	"log"
	"net"
)

type JoinRPCArgs struct {
	ID NodeAddress
}

type JoinRPCReply struct {
	Ok bool
}

func (node *Node) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
	log.Printf("node with ID: %v is joining the ring through: %v\n", args.ID, node.ID)
	reply.Ok = true
	return nil
}

func SendJoin(ID NodeAddress, addr *net.TCPAddr) error {
	args := JoinRPCArgs{}
	args.ID = ID
	reply := JoinRPCReply{}
	log.Printf("Joining ring at %v with ID %v\n", addr, ID)
	return call("Node.Join", addr, &args, &reply)
}
