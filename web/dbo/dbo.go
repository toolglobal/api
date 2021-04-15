package dbo

import (
	"github.com/wolot/api/datamanager"
)

type DBO struct {
	dataM *datamanager.DataManager
}

func New(dataM *datamanager.DataManager) *DBO {
	return &DBO{
		dataM: dataM,
	}
}
