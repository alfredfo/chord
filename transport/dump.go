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
  reply.DumpString = fmt.Sprintf("tcpAdd: %v,\n succ: %v,\n pred: %v,\n kv stored: %v,\n backup: %v,\n fingers: %v\n",
		tp.Node.NodeInfo, tp.Node.Successors[0], tp.Node.Predecessor, tp.Node.Bucket, tp.Node.Backup, tp.Node.FingerTable)
	return nil
}

func SendDump(addr *net.TCPAddr) (string, error) {
	args := DumpRPCArgs{}
	reply := DumpRPCReply{}
	err := call("TransportNode.Dump", addr, &args, &reply)

	return reply.DumpString, err
}
