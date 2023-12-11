package transport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"log"
	"net"

	"github.com/alfredfo/chord/api"
	"github.com/alfredfo/chord/hashing"
)

type SetRPCArgs struct {
	Key   api.Key
	Value api.Value
}
type SetRPCReply struct{}

type GetRPCArgs struct {
	Key api.Key
}
type GetRPCReply struct {
	Value api.Value
}

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func encryptIt(value []byte, keyPhrase string) ([]byte, error) {

	aesBlock, err := aes.NewCipher([]byte(mdHashing(keyPhrase)))
	if err != nil {
		log.Printf("aesBlock: %v", err)
		return nil, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Printf("gcmInstance: %v", err)
		return nil, err
	}

	nonce := make([]byte, gcmInstance.NonceSize())

	return gcmInstance.Seal(nonce, nonce, value, nil), nil
}

func (tp *TransportNode) Set(args *SetRPCArgs, reply *SetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	log.Printf("set\n")
	tp.Node.Bucket[args.Key] = []byte(args.Value)
	log.Printf("current val in node %v bucket: %v", tp.Node.NodeInfo, tp.Node.Bucket)

	return nil
}

func decryptIt(ciphered []byte, keyPhrase string) ([]byte, error) {
	hashedPhrase := mdHashing(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		log.Printf("aesBlock: %v", err)
		return nil, err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Printf("gcmInstance: %v", err)
		return nil, err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Printf("originalText: %v", err)
		return nil, err
	}
	return originalText, nil
}

func (tp *TransportNode) Get(args *GetRPCArgs, reply *GetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()

	reply.Value = tp.Node.Bucket[args.Key]
	return nil
}

func (tp *TransportNode) Delete(args *GetRPCArgs, reply *GetRPCReply) error {
	tp.Node.Mu.Lock()
	defer tp.Node.Mu.Unlock()
	reply.Value = tp.Node.Bucket[args.Key]
	delete(tp.Node.Bucket, args.Key)
	log.Printf("current val in node %v bucket: %v", tp.Node.NodeInfo, tp.Node.Bucket)
	return nil
}

func SendSet(key api.Key, value api.Value, addr *net.TCPAddr, encriptKey string) error {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	n, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return err
	}
	args := SetRPCArgs{}
	log.Printf("kv: {%v, %v} will be stored at node: %v", key, value, n)
	args.Key = key
	value, err = encryptIt(value, encriptKey)
	if err != nil {
		log.Println(err)
	}

	args.Value = value
	reply := SetRPCReply{}
	log.Println("before call...")

	nsucc, err := SendAskSuccessors(&n.TCPAddr)

	if err == nil {
		err = call("TransportNode.Set", &nsucc[0].TCPAddr, &args, &reply)
		if err != nil {
			log.Printf("could not send data to n's successor, no redundancy! err: %v", err)
		}
	}

	return call("TransportNode.Set", &n.TCPAddr, &args, &reply)

}

func SendGet(key api.Key, addr *net.TCPAddr, decriptKey string) (api.Value, error) {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	succ, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	args := GetRPCArgs{}
	args.Key = key
	reply := GetRPCReply{}
	err = call("TransportNode.Get", &succ.TCPAddr, &args, &reply)
	val, err := decryptIt(reply.Value, decriptKey)
	if err != nil {
		log.Println(err)
	}
	return val, err
}

func SendDelete(key api.Key, addr *net.TCPAddr) (api.Value, error) {
	keyHash := hashing.HashStringToBigInt(key).String()
	log.Printf("Hash value for key :%v is %v", key, keyHash)
	succ, err := SendFindSuccessor(keyHash, addr)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	args := GetRPCArgs{}
	args.Key = key
	reply := GetRPCReply{}
	err = call("TransportNode.Delete", &succ.TCPAddr, &args, &reply)
	return reply.Value, err
}
