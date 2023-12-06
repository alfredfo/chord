package transport

import (
	"log"
	"net"
	"github.com/alfredfo/chord/api"
)

type NotifyRPCArgs struct {}

type NotifyRPCReply struct {
	SuccessorList []api.NodeAddress
}

func (tp *TransportNode) Notify(args *NotifyRPCArgs, reply *NotifyRPCReply) error {
	log.Println("recieved notify call")
	reply.SuccessorList = tp.Node.Successors
	return nil
}

func SendNotify(ID api.NodeAddress, addr *net.TCPAddr) error {
	args := NotifyRPCArgs{}
	reply := NotifyRPCReply{}
	log.Printf("Notifying ring at %v with ID %v\n", addr, ID)
	err := call("TransportNode.Notify", addr, &args, &reply)
	if err != nil {
		log.Printf("error sending Notify to %v: %v\n", ID, err)
	}
	log.Printf("sugma balls: %v\n", reply.SuccessorList)
	return nil
}
