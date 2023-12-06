package transport

import (
	"net"

	"github.com/alfredfo/chord/api"
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
	args.Key = key
	args.Value = value
	reply := SetRPCReply{}
	return call("TransportNode.Set", addr, &args, &reply)
}

func SendGet(key api.Key, addr *net.TCPAddr) (api.Value, error) {
	args := GetRPCArgs{}
	args.Key = key
	reply := GetRPCReply{}
	err := call("TransportNode.Get", addr, &args, &reply)
	return reply.Value, err
}
