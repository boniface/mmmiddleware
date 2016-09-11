package dbconn

import (
	"testing"
	"github.com/pborman/uuid"
	"fmt"
)

func TestCreateKeyspace(t *testing.T) {
	CreateKeyspace()
}
func TestInsertRowTable(t *testing.T){
	type Itest struct {
		Year string
		Month string
		Day string
		Id string
		OrgCode string
		Code string
	}
	itest :=Itest{}
	itest.Code="1000/0001"
	itest.Day = "04"
	itest.Id = uuid.New()
	itest.OrgCode = "MM001"
	itest.Month = "09"
	itest.Year = "2016"
	InsertRowTable(itest,"testInsert")
}
func TestRunQueryCassWithFeedback(t *testing.T){
	qry :="select * from testInsert"
	str :=RunQueryCassCollection(qry)
	fmt.Println("===>>? ",str)
}