package transport

import (
	"net"
)

type CheckPredecessorRPCArgs struct {}

type CheckPredecessorRPCReply struct {}

func SendCheckPredecessor(addr *net.TCPAddr) error {
	args := CheckPredecessorRPCArgs{}
	reply := CheckPredecessorRPCReply{}
	
	err := call("TransportNode.CheckPredecessor", addr, &args, &reply)
	return err
}

func (tp *TransportNode) CheckPredecessor(args *AskPredecessorRPCArgs, reply *AskPredecessorRPCReply) error { 
	return nil
}
