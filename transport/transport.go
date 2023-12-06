package transport

import (
	"github.com/alfredfo/chord/api"
)

type Transport interface {
	Join(ID api.NodeAddress)
}

type TransportNode struct {
	Node *api.Node
	get *Transport
}
