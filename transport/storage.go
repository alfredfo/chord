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
	tp.Node.Bucket[args.Key] = args.Value
	log.Printf("current val in node %v bucket: %v", tp.Node.ID, tp.Node.Bucket)
	return nil
}

func (tp *TransportNode) Get(args *GetRPCArgs, reply *GetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	reply.Value = tp.Node.Bucket[args.Key]
	return nil
}

func SendSet(key api.Key, value api.Value, addr *net.TCPAddr) error {
	args := SetRPCArgs{}
	keyHash := hashing.HashString(key)
	args.Key = keyHash.String()
	args.Value = value
	reply := SetRPCReply{}
	return call("TransportNode.Set", addr, &args, &reply)
}

func SendGet(key api.Key, addr *net.TCPAddr) (api.Value, error) {
	args := GetRPCArgs{}
	keyHash := hashing.HashString(key)
	args.Key = keyHash.String()
	reply := GetRPCReply{}
	err := call("TransportNode.Get", addr, &args, &reply)
	return reply.Value, err
}
