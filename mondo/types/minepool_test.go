package types

import (
	"fmt"
	"testing"
	"time"
)

func TestMinePool_encode_decode(t *testing.T) {
	var src MinePool
	src.GenesisTime = time.Now()
	src.LastBlockHeight = 1024
	src.LastBlockTime = time.Now().Add(time.Second * 12)

	b := src.ToBytes()

	var dst MinePool
	dst.FromBytes(b)

	fmt.Println(src)

}
