package main

import (
	"crypto/sha1"
	"math/big"
	"github.com/alfredfo/chord/api"
	"net"
)

//const keySize = sha1.Size * 8
const keySize = 4
var two = big.NewInt(2)
var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

func jump(node *api.Node, fingerentry int) *api.NodeAddress {
	n := big.Int(node.ID)

	fingerentryminus1 := big.NewInt(int64(fingerentry) - 1)
	jump := new(big.Int).Exp(two, fingerentryminus1, nil)
	sum := new(big.Int).Add(&n, jump)

	return new(big.Int).Mod(sum, hashMod)
}

func hashString(elt string) *api.NodeAddress {
    hasher := sha1.New()
    hasher.Write([]byte(elt))
    return new(big.Int).SetBytes(hasher.Sum(nil))
}

func hashAddress(addr *net.TCPAddr) *api.NodeAddress {
	return hashString(addr.String())
}