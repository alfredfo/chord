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
	log.Printf("node with ID: %v is joining the ring through: \n", args.ID, node.ID)
	reply.Ok = true
	return nil
}

func (node *Node) SendJoin(addr *net.TCPAddr) error {
	args := JoinRPCArgs{}
	args.ID = node.ID
	reply := JoinRPCReply{}
	return call("Node.Join", addr, &args, &reply)
}
