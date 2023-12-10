package transport

import (
	"fmt"
	"net"
	// "github.com/alfredfo/chord/api"
)

type DumpRPCArgs struct {
}
type DumpRPCReply struct {
	DumpString string
}

func (tp *TransportNode) Dump(args *DumpRPCArgs, reply *DumpRPCReply) error {
	reply.DumpString = fmt.Sprintf("tcpAdd: %v, succ: %v, pred: %v, kv stored: %v, successors: %v\n",
		tp.Node.NodeInfo, tp.Node.Successors[0], tp.Node.Predecessor, tp.Node.Bucket, tp.Node.Successors)
	return nil
}

func SendDump(addr *net.TCPAddr) (string, error) {
	args := DumpRPCArgs{}
	reply := DumpRPCReply{}
	err := call("TransportNode.Dump", addr, &args, &reply)

	return reply.DumpString, err
}
