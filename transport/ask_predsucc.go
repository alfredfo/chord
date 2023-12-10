package transport

import (
	"net"

	"github.com/alfredfo/chord/api"
)

type AskPredecessorRPCArgs struct{}

type AskPredecessorRPCReply struct {
	Predecessor api.NodeInfoType
}

func SendAskPredecessor(addr *net.TCPAddr) (api.NodeInfoType, error) {
	args := AskPredecessorRPCArgs{}
	reply := AskPredecessorRPCReply{}

	err := call("TransportNode.AskPredecessor", addr, &args, &reply)
	return reply.Predecessor, err
}

func (tp *TransportNode) AskPredecessor(args *AskPredecessorRPCArgs, reply *AskPredecessorRPCReply) error {
	reply.Predecessor = tp.Node.Predecessor
	return nil
}

type AskSuccessorsRPCArgs struct{}

type AskSuccessorsRPCReply struct {
	Successors []api.NodeInfoType
}

func SendAskSuccessors(addr *net.TCPAddr) ([]api.NodeInfoType, error) {
	args := AskSuccessorsRPCArgs{}
	reply := AskSuccessorsRPCReply{}

	err := call("TransportNode.AskSuccessors", addr, &args, &reply)
	return reply.Successors, err
}

func (tp *TransportNode) AskSuccessors(args *AskSuccessorsRPCArgs, reply *AskSuccessorsRPCReply) error {
	reply.Successors = tp.Node.Successors
	return nil
}
