package transport

import (
	"github.com/alfredfo/chord/api"
)

type Transport interface {
	Join(args *JoinRPCArgs, reply *JoinRPCReply)
	Notify(args *NotifyRPCArgs, reply *NotifyRPCReply)
	Set(args *SetRPCArgs, reply *SetRPCReply)
	Get(args *GetRPCArgs, reply *GetRPCReply)
}

type TransportNode struct {
	Node *api.Node
	get  *Transport
}
