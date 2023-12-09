package transport

import (
	"log"
	// "math/big"
	"net"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
)

type JoinRPCArgs struct {
	// ID api.NodeAddress
	ToJoinNode api.NodeInfoType
}

type JoinRPCReply struct {
	Successor  api.NodeInfoType
	Ok         bool
}

// Join only deal with join logic, the tp.Node must be args.NodeInfo' successor
// We use find_sucessor to find the sucessor first, then call this
func (tp *TransportNode) Join(args *JoinRPCArgs, reply *JoinRPCReply) error {
	log.Printf("Node %v is trying to join the ring through: %v\n", args.ToJoinNode, tp.Node.NodeInfo)
	reply.Ok = true
	//tp.Node.Mu.Lock()
	//defer tp.Node.Mu.Unlock()

	log.Printf("suc %v\n", tp.Node.Successor)
	r, _ := SendFindSuccessor(args.ToJoinNode.ID, &tp.Node.Successor.TCPAddr)
	reply.Successor = r
	
	log.Printf("Join: %v\n", tp.Node.NodeInfo)
	log.Printf("Join: %v\n", tp.Node.Successor)
	return nil
}

// node will join ring through addr
// return the sucessor and predcessor of node after join
func SendJoin(node *api.Node, addr *net.TCPAddr) (api.NodeInfoType, error) {
	args := JoinRPCArgs{}
	args.ToJoinNode = node.NodeInfo
	reply := JoinRPCReply{}

	err := call("TransportNode.Join", addr, &args, &reply)
	return reply.Successor, err
}

type CopyDataRPCArgs struct {
  NodeInfo api.NodeInfoType
}

type CopyDataRPCReply struct {
  Data  map[api.Key]api.Value 
}

func (tp *TransportNode) CopyData(args *CopyDataRPCArgs, reply *CopyDataRPCReply) error {
	log.Printf("Node %v is copying data to %v\n", tp.Node.NodeInfo, args.NodeInfo)
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()

  reply.Data = make(map[api.Key]api.Value)
  for k, v := range tp.Node.Bucket {
    keyHash := hashing.HashStringToBigInt(k).String()
    log.Printf("==== %v %v %v\n", tp.Node.NodeInfo.ID, keyHash, args.NodeInfo.ID)
    if hashing.SBetween(tp.Node.NodeInfo.ID, keyHash, args.NodeInfo.ID, true) {
      reply.Data[k] = v
    } 
  }

	return nil
}

func SendCopyData(nodeInfo api.NodeInfoType, addr *net.TCPAddr) (map[api.Key]api.Value, error) {
	args := CopyDataRPCArgs{}
	args.NodeInfo = nodeInfo
	reply := CopyDataRPCReply{}

	err := call("TransportNode.CopyData", addr, &args, &reply)
	return reply.Data, err
}


type StoreDataRPCArgs struct {
  Data map[api.Key]api.Value
}

type StoreDataRPCReply struct {
}

func (tp *TransportNode) StoreData(args *StoreDataRPCArgs, reply *StoreDataRPCReply) error {
	log.Printf("Node %v is storing data", tp.Node.NodeInfo)
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()

  for k, v := range args.Data {
    tp.Node.Bucket[k] = v
  }

	return nil
}

func SendStoreData(kvmap map[api.Key]api.Value, addr *net.TCPAddr) error {
	args := StoreDataRPCArgs{}
	args.Data = kvmap
	reply := CopyDataRPCReply{}
	err := call("TransportNode.StoreData", addr, &args, &reply)
	return err
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
func SendChangeSucessor(nodeInfo api.NodeInfoType, addr *net.TCPAddr) error {
	args := ChangeSuccessorRPCArgs{}
	args.NewSuccessor = nodeInfo
	reply := ChangeSuccessorRPCReply{}

	err := call("TransportNode.ChangeSucessor", addr, &args, &reply)
	return err
}
