package transport

import (
	"github.com/alfredfo/chord/api"
)

type Transport interface {
	Join(args *JoinRPCArgs, reply *JoinRPCReply)
	Notify(args *NotifyRPCArgs, reply *NotifyRPCReply)
}

type TransportNode struct {
	Node *api.Node
	get *Transport
}
