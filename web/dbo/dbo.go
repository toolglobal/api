package dbo

import (
	"github.com/toolglobal/api/datamanager"
)

type DBO struct {
	dataM *datamanager.DataManager
}

func New(dataM *datamanager.DataManager) *DBO {
	return &DBO{
		dataM: dataM,
	}
}
