package api

import (
	"log"
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

func (node *Node) FindSuccessor(*int) *Node {
	return nil
}

func (node *Node) Notify(childNode *Node) {

}
