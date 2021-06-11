package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func Test_Config(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	cfg := Params{
		DPOSBeginHeight:  1,
		DPOSEachHeight:   2,
		DPOSMaxNodeNum:   3,
		NodeWorkMortgage: 4,
		NodeMinMortgage:  5,
		NodeMinCollect:   6,
		UserMinMortgage:  7,
		UserMinCollect:   8,
	}
	SaveMondoParams(db, &cfg)

	ncfg := LoadMondoParams(db)

	fmt.Println(cfg)
	fmt.Println(ncfg)
	if !reflect.DeepEqual(cfg, *ncfg) {
		t.Fatal("not equal")
	}
}
