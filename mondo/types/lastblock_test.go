package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func Test_marshal(t *testing.T) {
	var block LastBlockInfo
	block.Height = 123
	block.StateRoot = common.HexToHash("0x12222")

	block.ValidatorUpdates = append(block.ValidatorUpdates, ABCIValidator{
		PubKey: [32]byte{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 233, 34},
		Power:  10,
	})

	block.ValidatorUpdates = append(block.ValidatorUpdates, ABCIValidator{
		PubKey: [32]byte{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 233, 35},
		Power:  11,
	})

	b := block.ToBytes()
	var nb LastBlockInfo
	nb.FromBytes(b)
	fmt.Println(nb.Height)

	for _, v := range nb.ValidatorUpdates {
		fmt.Println(v.PubKey)
	}

}
