package hashing

import (
	"crypto/sha1"
	"github.com/alfredfo/chord/api"
	"log"
	"math/big"
	"net"
)

var two = big.NewInt(2)
var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(api.KeySize), nil)

func Jump(ID api.NodeID, fingerentry int) api.NodeID {
	n := NodeIDToBigInt(ID)

	//	log.Printf("jump | n: %v", n)
	fingerentryminus1 := big.NewInt(int64(fingerentry) - 1)

	jump := new(big.Int).Exp(two, fingerentryminus1, nil)
	//	log.Printf("jump | jump: %v", jump)
	sum := new(big.Int).Add(n, jump)
	//	log.Printf("jump | sum: %v", sum)
	return new(big.Int).Mod(sum, hashMod).String()
}

func Between(start, elt, end *big.Int, inclusive bool) bool {
	if end.Cmp(start) > 0 {
		return (start.Cmp(elt) < 0 && elt.Cmp(end) < 0) || (inclusive && elt.Cmp(end) == 0)
	} else {
		return start.Cmp(elt) < 0 || elt.Cmp(end) < 0 || (inclusive && elt.Cmp(end) == 0)
	}
}

func NodeIDToBigInt(ID api.NodeID) *big.Int {
	n := new(big.Int)
	n, ok := n.SetString(ID, 10)
	if ok == false {
		log.Printf("NodeIDToBigInt: %v", ID)
	}
	return n
}

func SBetween(sstart, selt, send api.NodeID, inclusive bool) bool {
	start := NodeIDToBigInt(sstart)
	elt := NodeIDToBigInt(selt)
	end := NodeIDToBigInt(send)

	if end.Cmp(start) > 0 {
		return (start.Cmp(elt) < 0 && elt.Cmp(end) < 0) || (inclusive && elt.Cmp(end) == 0)
	} else {
		return start.Cmp(elt) < 0 || elt.Cmp(end) < 0 || (inclusive && elt.Cmp(end) == 0)
	}
}

func HashStringToBigInt(elt string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(elt))
	hash := new(big.Int).SetBytes(hasher.Sum(nil))
	hash = new(big.Int).Mod(hash, hashMod)
	return hash
}

func HashTcpAddressToString(addr *net.TCPAddr) string {
	return HashStringToBigInt(addr.String()).String()
}
