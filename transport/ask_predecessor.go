package transport

import (
	"net"

	"github.com/alfredfo/chord/api"
)

type AskPredecessorRPCArgs struct{}

type AskPredecessorRPCReply struct {
	Successors  []api.NodeInfoType
	Predecessor api.NodeInfoType
}

func SendAskPredecessor(addr *net.TCPAddr) ([]api.NodeInfoType, api.NodeInfoType, error) {
	args := AskPredecessorRPCArgs{}
	reply := AskPredecessorRPCReply{}

	err := call("TransportNode.AskPredecessor", addr, &args, &reply)
	return reply.Successors, reply.Predecessor, err
}

func (tp *TransportNode) AskPredecessor(args *AskPredecessorRPCArgs, reply *AskPredecessorRPCReply) error {
	reply.Predecessor = tp.Node.Predecessor
	succSlice := make([]api.NodeInfoType, 0)
	succSlice = append(succSlice, tp.Node.Successor)
	reply.Successors = succSlice
	return nil
}
