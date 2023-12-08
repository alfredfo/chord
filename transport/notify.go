package transport

import (
	"log"
	// "net"

	"github.com/alfredfo/chord/api"
)

type NotifyRPCArgs struct{}

type NotifyRPCReply struct {
	Successor api.NodeInfoType
}

func (tp *TransportNode) Notify(args *NotifyRPCArgs, reply *NotifyRPCReply) error {
	log.Println("recieved notify call")
	reply.Successor = tp.Node.Successor
	return nil
}

func SendNotify(node *api.Node) error {
	args := NotifyRPCArgs{}
	reply := NotifyRPCReply{}
	log.Printf("Notifying ring at %v with ID %v\n", node.NodeInfo, node.NodeInfo)
	err := call("TransportNode.Notify", &node.NodeInfo.TCPAddr, &args, &reply)
	if err != nil {
		log.Printf("error sending Notify to %v: %v\n", node.NodeInfo, err)
	}
	log.Printf("sugma balls: %v\n", reply.Successor)
	return nil
}
