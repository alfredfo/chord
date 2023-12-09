package transport

import (
	"log"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
)

type NotifyRPCArgs struct{
	NodeInfo api.NodeInfoType
}

type NotifyRPCReply struct {}

func (tp *TransportNode) Notify(args *NotifyRPCArgs, reply *NotifyRPCReply) error {
	// log.Println("recieved notify call")

	ni := args.NodeInfo
	n := ni.ID
	pred := tp.Node.Predecessor
	
	if pred.ID == "" || hashing.SBetween(pred.ID, n, tp.Node.NodeInfo.ID, false) {
		tp.Node.Predecessor = ni
	}
	
	return nil
}

func SendNotify(node *api.Node, other api.NodeInfoType) error {
	args := NotifyRPCArgs{}
	reply := NotifyRPCReply{}
	// log.Printf("Notifying ring at %v with ID %v\n", node.NodeInfo, node.NodeInfo)

	args.NodeInfo = node.NodeInfo
	
	err := call("TransportNode.Notify", &node.Successor.TCPAddr, &args, &reply)
	if err != nil {
		log.Printf("error sending Notify to %v: %v\n", node.NodeInfo, err)
	}
	// log.Printf("sugma balls\n")
	return nil
}
