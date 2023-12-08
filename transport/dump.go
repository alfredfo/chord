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
  reply.DumpString = fmt.Sprintf("tcpAdd: %v, succ: %v, pred: %v, kv stored: %v\n", 
    tp.Node.Address, tp.Node.Successor, tp.Node.Predecessor, tp.Node.Bucket)
	return nil
}

func SendDump(addr *net.TCPAddr) (string, error) {
	args := DumpRPCArgs{}
	reply := DumpRPCReply{}
	err := call("TransportNode.Dump", addr, &args, &reply)

	return reply.DumpString, err
}




