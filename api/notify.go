package api

type NotifyRPCArgs struct {
	ID NodeAddress
}

type NotifyRPCReply struct {
	Ok bool
}

func (node *Node) Notify(childNode *Node) {

}
