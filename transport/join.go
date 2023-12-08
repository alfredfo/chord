package transport

import (
	"log"
	"math/big"
	"net"

	"github.com/alfredfo/chord/api"
)

type JoinRPCArgs struct {
	ID api.NodeAddress
}

type JoinRPCReply struct {
	FoundSuccessor api.NodeAddress
	Ok             bool
}

func (tp *TransportNode) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
	log.Printf("node with ID: %v is joining the ring through: %v\n", args.ID.String(), tp.Node.ID.String())
	reply.Ok = true
	tp.Node.Successor = nil
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	// SendFindSuccessor(tp.Node.ID, nil)
	//tp.Node.FingerTable = append(tp.Node.FingerTable, args.ID)
	// s := make(map[string]int)
	// tp.Node.Successor = map[]
	return nil
}

func SendJoin(ID api.NodeAddress, addr *net.TCPAddr) (api.NodeAddress, error) {
	args := JoinRPCArgs{}
	args.ID = (big.Int)(ID)
	reply := JoinRPCReply{}
	log.Printf("Trying to join ring through %v with ID %v\n", addr, args.ID.String())
	err := call("TransportNode.Join", addr, &args, &reply)
	return reply.FoundSuccessor, err
}
