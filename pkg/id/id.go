// pkg/id/id.go
package id

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
)

type NodeID [20]byte
type RPCID  [20]byte

func NewRandomRPC() RPCID {
	var r RPCID
	_, _ = rand.Read(r[:])
	return r
}

func FromHex(s string) (NodeID, error) {
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 20 {
		return NodeID{}, err
	}
	var n NodeID
	copy(n[:], b)
	return n, nil
}

func ToHex(n [20]byte) string { return hex.EncodeToString(n[:]) }

func HashSHA1(data []byte) NodeID {
	sum := sha1.Sum(data)
	var n NodeID
	copy(n[:], sum[:])
	return n
}

func XORDistance(a, b NodeID) [20]byte {
	var d [20]byte
	for i := 0; i < 20; i++ {
		d[i] = a[i] ^ b[i]
	}
	return d
}