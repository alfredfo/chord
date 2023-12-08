package transport

import (
	"net"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
	"log"
)

type FindSuccessorRPCArgs struct {
	ID api.NodeID
}

type FindSuccessorRPCReply struct {
	Successor api.NodeInfoType
}

func SendFindSuccessor(ID api.NodeID, addr *net.TCPAddr) (api.NodeInfoType, error) {
	args := FindSuccessorRPCArgs{}
	reply := FindSuccessorRPCReply{}
	args.ID = ID
	
	err := call("TransportNode.FindSuccessor", addr, &args, &reply)
	return reply.Successor, err
}

func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
	ID := args.ID
	ourID := tp.Node.NodeInfo.ID
	succ := tp.Node.Successor
	log.Printf("lel %v\n", succ)
	succID := succ.ID

	if ourID == succID {
		reply.Successor = tp.Node.NodeInfo
		return nil
	}

	log.Printf("sugma %v | %v | %v\n", ID, ourID, succID)
	if hashing.SBetween(ourID, ID, succID, true) {
		reply.Successor = succ
	} else {
		nodeInfo, err := SendFindSuccessor(ID, &tp.Node.Successor.TCPAddr)
		if err != nil {
			return err
		}
		reply.Successor = nodeInfo
	}

	return nil
}



// n.find successor(id)
// if (id ∈ (n,successor])
// return successor;
// else
// n
// 0 = closest preceding node(id);
// return n
// 0
// .find successor(id);

// func (tp *TransportNode) FindSuccessor(args *FindSuccessorRPCArgs, reply *FindSuccessorRPCReply) error {
//   var (
//     joinID int
//     succID int
//     predID int
//     err error
//   )
//   // base case, when the cursive call returns
//   for k, _ := range args.ToJoinNode {
//     log.Printf("join k: %v", k)
//     joinID, err = strconv.Atoi(k)
//     if err != nil {
//       log.Println(err)
//     }
//   }
//
//   for k, _ := range tp.Node.Successor {
//     succID, err = strconv.Atoi(k)
//     if err != nil {
//       log.Println(err)
//
//     }
//   }
//
//   for k, _ := range tp.Node.Predecessor {
//     predID, err = strconv.Atoi(k)
//     if err != nil {
//       log.Println(err)
//
//     }
//   }
//
//
//   currID, err := strconv.Atoi(tp.Node.ID.String())
//   if err != nil {
//     log.Println(err)
//   }
//   log.Printf("[find_sucessor] asking %v and it's succ: %v what should be: %v's sucessor", currID, succID, joinID)
//
//   // only one Node
//   if succID == currID {
//     reply.Successor = tp.Node.Successor
//     return nil
//   }
//
//   // TODO when to terminate recursive call
//   // max_relay := 10 // max time to ask successor
//   // if joinID > currID && joinID <= succID {
//   //   log.Printf("successor found!, %v is the sucessor of %v", args.ToJoinNode, tp.Node.Successor)
//   //   reply.Successor = tp.Node.Successor
//   //   return nil
//   // } else {
//   //   curr := tp.Node.Successor
//   //   // recursive call
//   //   for _, addr := range curr {
//   //     err := call("TransportNode.FindSuccessor", &addr, args, reply)
//   //     if err != nil {
//   //       log.Println(err)
//   //     }
//   //   }
//   // }
//   if joinID > currID {
//     if joinID <= succID {
//       log.Printf("successor found!, %v is the sucessor of %v", args.ToJoinNode, tp.Node.Successor)
//       reply.Successor = tp.Node.Successor
//       return nil
//     } else {
//       curr := tp.Node.Successor
//       // recursive call
//       for _, addr := range curr {
//         err := call("TransportNode.FindSuccessor", &addr, args, reply)
//         if err != nil {
//           log.Println(err)
//         }
//       }
//     }
//   } else { // joinID < currID
//     if joinID > predID {
//       log.Printf("successor found!, %v is the sucessor of %v", args.ToJoinNode, tp.Node.Predecessor)
//         reply.Successor = tp.Node.Predecessor
//     } else { // joinID <= predID
//       curr := tp.Node.Predecessor
//       // recursive call
//       for _, addr := range curr {
//         err := call("TransportNode.FindSuccessor", &addr, args, reply)
//         if err != nil {
//           log.Println(err)
//         }
//       }
//     }
//   }
//
// 	return nil
// }
//

//
//
//
//
