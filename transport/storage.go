package transport

import (
	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"log"
	"net"
)

type SetRPCArgs struct {
	Key   api.Key
	Value api.Value
}
type SetRPCReply struct{}

type GetRPCArgs struct {
	Key api.Key
}
type GetRPCReply struct {
	Value api.Value
}

func (tp *TransportNode) Set(args *SetRPCArgs, reply *SetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	log.Printf("set\n")
	tp.Node.Bucket[args.Key] = args.Value
	log.Printf("current val in node %v bucket: %v", tp.Node.NodeInfo, tp.Node.Bucket)
	return nil
}

func (tp *TransportNode) Get(args *GetRPCArgs, reply *GetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	reply.Value = tp.Node.Bucket[args.Key]
	return nil
}

func (tp *TransportNode) Delete(args *GetRPCArgs, reply *GetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	reply.Value = tp.Node.Bucket[args.Key]
	delete(tp.Node.Bucket, args.Key)
	log.Printf("current val in node %v bucket: %v", tp.Node.NodeInfo, tp.Node.Bucket)
	return nil
}

func SendSet(key api.Key, value api.Value, addr *net.TCPAddr) error {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	succ, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return err
	}
	args := SetRPCArgs{}
	log.Printf("kv: {%v, %v} will be stored at node: %v", key, value, succ)
	args.Key = key
	args.Value = value
	reply := SetRPCReply{}
	log.Println("before call...")
	return call("TransportNode.Set", &succ.TCPAddr, &args, &reply)
}

func SendGet(key api.Key, addr *net.TCPAddr) (api.Value, error) {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	succ, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return "", err
	}
	args := GetRPCArgs{}
	args.Key = key
	reply := GetRPCReply{}
	err = call("TransportNode.Get", &succ.TCPAddr, &args, &reply)
	return reply.Value, err
}

func SendDelete(key api.Key, addr *net.TCPAddr) (api.Value, error) {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	succ, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return "", err
	}

	args := GetRPCArgs{}
	args.Key = key
	reply := GetRPCReply{}
	err = call("TransportNode.Delete", &succ.TCPAddr, &args, &reply)
	return reply.Value, err
}
