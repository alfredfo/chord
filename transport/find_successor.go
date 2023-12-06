package transport

import (
	"github.com/alfredfo/chord/api"
)

type FindSuccessorRPCArgs struct {
	ID api.NodeAddress
}
type FindSuccessorRPCReply struct {
	Successor api.NodeAddress
}

func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
	return nil
}

// func SendFindSuccessor(ID api.NodeAddress, addr *net.TCPAddr) (api.NodeAddress, error) {
// 	args := FindSuccessorRPCArgs{}
// 	reply := FindSuccessorRPCReply{}
//  successor, err := call("TransportNode.FindSuccessor", addr, &args, &reply)
// 	return successor, err
// }
