package transport

import (
	"log"
	// "math/big"
	"net"

	"github.com/alfredfo/chord/api"
)

type JoinRPCArgs struct {
	// ID api.NodeAddress
	ToJoinNode api.NodeInfoType
}

type JoinRPCReply struct {
	Successor  api.NodeInfoType
	Predcessor api.NodeInfoType
	Ok         bool
}

// Join only deal with join logic, the tp.Node must be args.NodeInfo' successor
// We use find_sucessor to find the sucessor first, then call this
func (tp *TransportNode) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
	log.Printf("Node %v is trying to join the ring through: %v\n", args.ToJoinNode, tp.Node.NodeInfo)
	reply.Ok = true
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()

	log.Printf("suc %v\n", tp.Node.Successor)
	reply.Predcessor = api.NodeInfoType{}
	r, _ := SendFindSuccessor(args.ToJoinNode.ID, &tp.Node.Successor.TCPAddr)
	reply.Successor = r
	
	log.Printf("Join: %v\n", tp.Node.NodeInfo)
	log.Printf("Join: %v\n", tp.Node.Successor)
	log.Printf("Join: %v\n", tp.Node.Predecessor)
	return nil
}

// node will join ring through addr
// return the sucessor and predcessor of node after join
func SendJoin(node *api.Node, addr *net.TCPAddr) (api.NodeInfoType, api.NodeInfoType, error) {
	args := JoinRPCArgs{}
	args.ToJoinNode = node.NodeInfo
	reply := JoinRPCReply{}

	err := call("TransportNode.Join", addr, &args, &reply)
	return reply.Successor, reply.Predcessor, err
}

type ChangeSuccessorRPCArgs struct {
	NewSuccessor api.NodeInfoType
}
type ChangeSuccessorRPCReply struct {
	
}

func (tp *TransportNode) ChangeSucessor(args *ChangeSuccessorRPCArgs, reply *ChangeSuccessorRPCReply) error {
	log.Printf("Node %v is changing sucessor to %v\n", tp.Node.NodeInfo, args.NewSuccessor)
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()

	tp.Node.Successor = args.NewSuccessor

	return nil
}

// change sucessor of addr to node
func SendChangeSucessor(node *api.Node, addr *net.TCPAddr) error {
	args := ChangeSuccessorRPCArgs{}
	args.NewSuccessor = node.NodeInfo
	reply := ChangeSuccessorRPCReply{}

	err := call("TransportNode.ChangeSucessor", addr, &args, &reply)
	return err
}
