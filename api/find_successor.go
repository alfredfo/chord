package api

type FindSuccessorRPCArgs struct {
	ID NodeAddress
}

type FindSuccessorRPCReply struct {
	Ok bool
}

func (node *Node) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) *Node {
	return nil
}
