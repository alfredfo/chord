package transport

import (
  "log"
	"net"

  // "github.com/alfredfo/chord/api"
)

type DumpRPCArgs struct {
}
type DumpRPCReply struct {
}

func (tp *TransportNode) Dump(args *DumpRPCArgs, reply *DumpRPCReply) error {
  log.Printf("tcpAdd: %v, succ: %v, pred: %v, kv stored: %v\n", 
    tp.Node.Address, tp.Node.Successor, tp.Node.Successor, tp.Node.Bucket)
	return nil
}

func SendDump(addr *net.TCPAddr) error {
	args := DumpRPCArgs{}
	reply := DumpRPCReply{}
	err := call("TransportNode.Dump", addr, &args, &reply)

	return err
}




