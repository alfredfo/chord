package transport

import (
	"net"
	"github.com/alfredfo/chord/api"
)

type FindSuccessorRPCArgs struct {}
type FindSuccessorRPCReply struct {}

func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
	return nil
}

func SendFindSuccessor(ID api.NodeAddress, addr *net.TCPAddr) error {
	args := FindSuccessorRPCArgs{}
	reply := FindSuccessorRPCReply{}

	return call("TransportNode.FindSuccessor", addr, &args, &reply)
}
