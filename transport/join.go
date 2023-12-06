package transport

import (
	"log"
	"net"

	"github.com/alfredfo/chord/api"
)

type JoinRPCArgs struct {
	ID api.NodeAddress
}

type JoinRPCReply struct {
	Ok bool
}

func (tp *TransportNode) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
	log.Printf("node with ID: %v is joining the ring through: %v\n", args.ID, tp.Node.ID)
	reply.Ok = true

	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	tp.Node.FingerTable = append(tp.Node.FingerTable, args.ID)
	// tp.Node.Predecessor = api.NodeAddress("")
	// SendFindSuccessor()
	// tp.Node.Successors = append()

	return nil
}

func SendJoin(ID api.NodeAddress, addr *net.TCPAddr) error {
	args := JoinRPCArgs{}
	args.ID = ID
	reply := JoinRPCReply{}
	log.Printf("Joining ring at %v with ID %v\n", addr, args.ID)
	return call("TransportNode.Join", addr, &args, &reply)
}
