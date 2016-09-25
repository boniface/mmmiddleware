package service

import (
	"testing"
	"fmt"
)

func TestRepFinanceStatement_LoadCustUploadData(t *testing.T) {
	rep :=RepFinanceStatement{}
	rep.OrgCode = "MM01"
	rep.Year = "2011"
	rep.Month = "3"
	rep.LoadCustUploadData()
	fmt.Println(" -->> ",len(rep.CustUploadData))
}
func TestRepFinanceStatement_RunRep(t *testing.T) {
	rep :=RepFinanceStatement{}
	rep.OrgCode = "MM01"
	rep.Year = "2011"
	rep.Month = "3"
	rep.LoadCustUploadData()
	fmt.Println(" -->> ",len(rep.CustUploadData))
	rep.RunRep()
}