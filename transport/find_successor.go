package transport

import (
	"github.com/alfredfo/chord/api"
	"net"
)

type FindSuccessorRPCArgs struct {
	ID api.NodeAddress
}
type FindSuccessorRPCReply struct {
	Successor api.NodeAddress
}

func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
	reply.Successor = api.NodeAddress{}
	return nil
}

func SendFindSuccessor(ID api.NodeAddress, addr *net.TCPAddr) (api.NodeAddress, error) {
	args := FindSuccessorRPCArgs{}
	reply := FindSuccessorRPCReply{}
	err := call("TransportNode.FindSuccessor", addr, &args, &reply)

	return reply.Successor, err
}




