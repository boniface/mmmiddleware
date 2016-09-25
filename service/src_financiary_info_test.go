package service

import (
	"testing"
	"fmt"
)

func TestFinanceInfo_LoadOneYeaUploadedInfo(t *testing.T) {
	info :=FinanceInfo{}
	info.Year = 2001
	info.Month = 1
	info.OrgCode = "MM01"
	info.ReqType = "year"
	info.LoadOneYeaUploadedInfo()

	rows :=info.UploadedInfo

	for _,row:=range rows{
		fmt.Println("==> Year:",row.Year," | month: ",row.Month)
	}
}