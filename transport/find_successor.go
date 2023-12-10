package transport

import (
	"math/big"
	"net"

	"log"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
)

type FindSuccessorRPCArgs struct {
	ID api.NodeID
}

type FindSuccessorRPCReply struct {
	Successor api.NodeInfoType
}

func SendFindSuccessor(ID api.NodeID, addr *net.TCPAddr) (api.NodeInfoType, error) {
	args := FindSuccessorRPCArgs{}
	reply := FindSuccessorRPCReply{}
	args.ID = ID

	err := call("TransportNode.FindSuccessor", addr, &args, &reply)
	return reply.Successor, err
}

func ClosestPrecedingNode(node *api.Node, ID api.NodeID) api.NodeInfoType {
	// TODO: actual value ?
	for i := 8; i > 0; i-- {
		finger := node.FingerTable[big.NewInt(int64(i)).String()]
		if finger.ID != "" {
			if hashing.SBetween(node.NodeInfo.ID, finger.ID, ID, false) {
				return finger
			}
		}
	}
	return node.NodeInfo
}

func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
	ID := args.ID
	ourID := tp.Node.NodeInfo.ID
	succ := tp.Node.Successors[0]

	succID := succ.ID

	if ourID == succID {
		reply.Successor = tp.Node.NodeInfo
		return nil
	}

	log.Printf("sugma %v | %v | %v\n", ID, ourID, succID)
	if hashing.SBetween(ourID, ID, succID, true) && succID != "" {
		reply.Successor = succ
	} else {
		// FingerTable method
		// closestNode := ClosestPrecedingNode(tp.Node, ID)
		// nodeInfo, err := SendFindSuccessor(ID, &closestNode.TCPAddr)

		nodeInfo, err := SendFindSuccessor(ID, &succ.TCPAddr)

		if err != nil {
			return err
		}

		reply.Successor = nodeInfo
	}

	return nil
}
