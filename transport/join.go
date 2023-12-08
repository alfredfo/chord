package transport

import (
	"log"
	// "math/big"
	"net"

	"github.com/alfredfo/chord/api"
)

type JoinRPCArgs struct {
	// ID api.NodeAddress
  NodeInfo map[string]net.TCPAddr
}

type JoinRPCReply struct {
	FoundSuccessor map[string]net.TCPAddr
  FoundPredcessor map[string]net.TCPAddr

	Ok             bool
}

func (tp *TransportNode) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
  for k, _ := range args.NodeInfo {
    log.Printf("node with ID: %v is joining the ring through: %v\n", k, tp.Node.ID.String())
  }
	reply.Ok = true
	// tp.Node.Successor = nil
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	// SendFindSuccessor(tp.Node.ID, nil)
	//tp.Node.FingerTable = append(tp.Node.FingerTable, args.ID)
	s := make(map[string]net.TCPAddr)
  s[tp.Node.ID.String()] = *tp.Node.Address
	// tp.Node.Successor = map[]
  reply.FoundSuccessor = s
  reply.FoundPredcessor = tp.Node.Predecessor
  
  tp.Node.Predecessor = args.NodeInfo
  
  // successor and node same = single node in ring ? 
  for k, v := range tp.Node.Successor {
    if k == tp.Node.ID.String() && v.String() == tp.Node.Address.String() {
      tp.Node.Successor = args.NodeInfo
    }
  }
	return nil
}

func SendJoin(node *api.Node, addr *net.TCPAddr) (map[string]net.TCPAddr, map[string]net.TCPAddr, error) {
	args := JoinRPCArgs{}
	// args.ID = (big.Int)(ID)
  s := make(map[string]net.TCPAddr)
  s[node.ID.String()] = *node.Address
  args.NodeInfo = s
	reply := JoinRPCReply{}
	log.Printf("Trying to join ring through %v with ID %v\n", addr, node.ID.String())
	err := call("TransportNode.Join", addr, &args, &reply)
	return reply.FoundSuccessor, reply.FoundPredcessor, err
}

func (tp *TransportNode) ChangeSucessor(args *JoinRPCArgs, reply *JoinRPCReply) error {
  for k, _ := range args.NodeInfo {
    log.Printf("change sucessor to ID: %v for node : %v\n", k, tp.Node.ID.String())
  }
	reply.Ok = true
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
  
  tp.Node.Successor = args.NodeInfo

	return nil
}

func SendChangeSucessor(node *api.Node, addr *net.TCPAddr) error {
  
  args := JoinRPCArgs{}
  s := make(map[string]net.TCPAddr)
  s[node.ID.String()] = *node.Address
  args.NodeInfo = s
	reply := JoinRPCReply{}
	
  err := call("TransportNode.ChangeSucessor", addr, &args, &reply)
  return err
}


