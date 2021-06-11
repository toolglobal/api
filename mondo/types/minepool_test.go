package types

import (
	"fmt"
	"testing"
	"time"
)

func TestDPosPool_encode_decode(t *testing.T) {
	var src DPosPool
	src.GenesisTime = time.Now()
	src.LastBlockHeight = 1024
	src.LastBlockTime = time.Now().Add(time.Second * 12)

	b := src.ToBytes()

	var dst DPosPool
	dst.FromBytes(b)

	fmt.Println(src)

}
