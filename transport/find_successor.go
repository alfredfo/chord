package transport

// import (
//   "log"
// 	"net"
// 	"strconv"
// 
// 	"github.com/alfredfo/chord/api"
// )
// 
// type FindSuccessorRPCArgs struct {
// 	ToJoinNode map[string]net.TCPAddr // node to join
//   BeginNode map[string]net.TCPAddr
// }
// type FindSuccessorRPCReply struct {
// 	Successor map[string]net.TCPAddr
// }

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
// func SendFindSuccessor(node *api.Node, addr *net.TCPAddr) (map[string]net.TCPAddr, error) {
// 	args := FindSuccessorRPCArgs{}
//   
//   nodeinfo := make(map[string]net.TCPAddr)
//   nodeinfo[node.ID.String()] = *node.Address
//   args.ToJoinNode = nodeinfo
// 	reply := FindSuccessorRPCReply{}
// 	err := call("TransportNode.FindSuccessor", addr, &args, &reply)
// 	return reply.Successor, err
// }
// 
// 
// 
// 
