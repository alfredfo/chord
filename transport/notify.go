package transport

import (
	"log"
	"github.com/alfredfo/chord/api"
)

type NotifyRPCArgs struct {}

type NotifyRPCReply struct {
	Successor api.NodeAddress
}

func (tp *TransportNode) Notify(args *NotifyRPCArgs, reply *NotifyRPCReply) error {
	log.Println("recieved notify call")
	reply.Successor = tp.Node.Successor.ID
	return nil
}

func SendNotify(node *api.Node) error {
	args := NotifyRPCArgs{}
	reply := NotifyRPCReply{}
	log.Printf("Notifying ring at %v with ID %v\n", node.Address, node.ID)
	err := call("TransportNode.Notify", node.Address, &args, &reply)
	if err != nil {
		log.Printf("error sending Notify to %v: %v\n", node.ID, err)
	}
	log.Printf("sugma balls: %v\n", reply.Successor)
	return nil
}
