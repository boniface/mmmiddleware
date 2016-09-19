package service

import (
	"hashcode.zm/mmmiddleware/model"
	"hashcode.zm/mmmiddleware/dbconn"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type IncomeItem struct {
	Desc string
	Cost float64
}

type IncomeStatement struct {
	OrgCode string
	Category string
	Total float64
	Month string
	Year string
	ListItem []IncomeItem
}


type RepFinanceStatement struct {
	OrgCode string
	Month string
	Year string
	TypeOfRep string // this can be Monthly , quarty , year
	IncomeStatement IncomeStatement
	Htable string
	Revenue IncomeStatement
	CostOfSale IncomeStatement
	GrossProfit IncomeStatement
	OtherRevenue IncomeStatement
	Expense IncomeStatement
	OperationProfit IncomeStatement
	Depreciation IncomeStatement
	InterestReceived IncomeStatement
	InterestPaid IncomeStatement
	ProfitBeforeIncomeTax IncomeStatement
	Tax IncomeStatement
	NetProfit IncomeStatement
	/* let get the target data */
	CustUploadData []model.CustomerUpload
}

type CustomerUpload struct {
	UpRef string
	OrgCode string
	Id string
	Refrence string
	Date string
	AccountingCode string
	Year int
	Month int
	Day int
	AccounttingSystem string
	TxnType string
	DebitValue string
	CreditValue string
	EntryDescription string
	EntryCategory string
	EntrySubCategory string
	CsvStringInput string
	MappingCode string
}


func(obj *RepFinanceStatement)RunRep(){
	//todo
	obj.LoadCustUploadData()
	obj.RevenueMethod()
	obj.CostOfSaleMethod()
	obj.OtherRevenueMethod()
	obj.GrossProfitMethod()
	obj.ExpenseMethod()
	obj.InterestPaidMethod()
	obj.InterestReceivedMethod()
	obj.DepreciationMethod()
}
func(obj *RepFinanceStatement)LoadCustUploadData(){
	qry :="select * from customerUpload where orgcode='"+obj.OrgCode+"' and year='"+obj.Year+"' "
	ls :=[]model.CustomerUpload{}
	rs :=dbconn.RunQueryCassCollection(qry)
	json.Unmarshal([]byte(rs),&ls)
	obj.CustUploadData = ls
}
/* component of finance */
func(obj *RepFinanceStatement)RevenueMethod(){
	statment :=obj.BuildData("SALES")
	obj.Revenue = statment
	fmt.Println("[SALES] --> ",statment.Total)
}
func(obj *RepFinanceStatement)CostOfSaleMethod(){
	statment :=obj.BuildData("COST OF SALES")
	obj.CostOfSale = statment
	fmt.Println("[COST OF SALES] --> ",statment.Total)
}
func(obj *RepFinanceStatement)GrossProfitMethod(){
	//todo
}
func(obj *RepFinanceStatement)OtherRevenueMethod(){
	statment :=obj.BuildData("OTHER INCOME")
	obj.OtherRevenue = statment
	fmt.Println("[OTHER INCOME] --> ",statment.Total)

}
func(obj *RepFinanceStatement)ExpenseMethod(){
	statment :=obj.BuildData("EXPENSES")
	obj.Expense = statment
	fmt.Println("[EXPENSES] --> ",statment.Total)
}
func(obj *RepFinanceStatement)OperationProfitMethod(){
	//todo
}
func(obj *RepFinanceStatement)DepreciationMethod(){
	//todo
}
func(obj *RepFinanceStatement)InterestReceivedMethod(){
	keysearch :="Interest Received"
	statement := obj.buildWithKeySearchDesc(keysearch)
	obj.InterestReceived = statement
	fmt.Println("[Interest Received] --> ",statement.Total,statement.ListItem)
}
func(obj *RepFinanceStatement)InterestPaidMethod(){
	keysearch :="Interest Paid"
	statement := obj.buildWithKeySearchDesc(keysearch)
	obj.InterestPaid = statement
	fmt.Println("[Interest Paid] --> ",statement.Total,statement.ListItem)



	/*searfor :=strings.ToLower("Interest Paid")
	qry :="select * from customerUpload where orgcode='"+obj.OrgCode+"' and year='"+obj.Year+"' "
	ls :=[]CustomerUpload{}
	rs :=dbconn.RunQueryCassCollection(qry)
	json.Unmarshal([]byte(rs),&ls)
	statement :=IncomeStatement{}
	statement.Month = obj.Month
	statement.Year = obj.Year
	statement.Category = "year"
	statement.Total = 0.0
	statement.OrgCode = obj.OrgCode
	for _,row:=range ls{
		orgvalue :=strings.ToLower(row.EntryDescription)
		if strings.Contains(orgvalue,searfor){
			item :=IncomeItem{}
			item.Desc = row.EntryDescription
			var cost float64 = 0

			if row.TxnType == "DEBIT"{
				cost ,_ =strconv.ParseFloat( row.DebitValue,64)
			}
			if row.TxnType == "CREDIT"{
				cost ,_ =strconv.ParseFloat( row.CreditValue,64)
				//cost = row.CreditValue
			}
			item.Cost = cost
			statement.Total = statement.Total+ cost
			statement.ListItem  = append(statement.ListItem,item)
		}
	}
	obj.InterestPaid = statement
*/
}

func(obj *RepFinanceStatement)TaxMethod(){
	//todo
}
func(obj *RepFinanceStatement)NetProfitMethod(){
	//todo
}


func(obj *RepFinanceStatement)ProfitBeforeIncomeTaxMethod(){
	//todo

}

func(obj *RepFinanceStatement)buildWithKeySearchDesc(keysearch string)IncomeStatement{
	//todo
	searfor :=strings.ToLower(keysearch)
	qry :="select * from customerUpload where orgcode='"+obj.OrgCode+"' and year='"+obj.Year+"' "
	ls :=[]CustomerUpload{}
	rs :=dbconn.RunQueryCassCollection(qry)
	json.Unmarshal([]byte(rs),&ls)
	statement :=IncomeStatement{}
	statement.Month = obj.Month
	statement.Year = obj.Year
	statement.Category = "year"
	statement.Total = 0.0
	statement.OrgCode = obj.OrgCode
	for _,row:=range ls{
		orgvalue :=strings.ToLower(row.EntryDescription)
		if strings.Contains(orgvalue,searfor){
			item :=IncomeItem{}
			item.Desc = row.EntryDescription
			var cost float64 = 0

			if row.TxnType == "DEBIT"{
				cost ,_ =strconv.ParseFloat( row.DebitValue,64)
			}
			if row.TxnType == "CREDIT"{
				cost ,_ =strconv.ParseFloat( row.CreditValue,64)
				//cost = row.CreditValue
			}
			item.Cost = cost
			statement.Total = statement.Total+ cost
			statement.ListItem  = append(statement.ListItem,item)
		}
	}
	return statement

}

func(obj *RepFinanceStatement)BuildData(cat string)IncomeStatement{
	qry :="select * from customerUpload where orgcode='"+obj.OrgCode+"' and year='"+obj.Year+"' "
	//ls2 :=[]model.CustomerUpload{}
	ls :=[]CustomerUpload{}
	rs :=dbconn.RunQueryCassCollection(qry)
	json.Unmarshal([]byte(rs),&ls)
	statement :=IncomeStatement{}
	statement.Month = obj.Month
	statement.Year = obj.Year
	statement.Category = "year"
	statement.Total = 0.0
	statement.OrgCode = obj.OrgCode
	for _,row:=range ls{
		if row.EntryCategory==cat{
			item :=IncomeItem{}
			item.Desc = row.EntryDescription
			var cost float64 = 0

			if row.TxnType == "DEBIT"{
				cost ,_ =strconv.ParseFloat( row.DebitValue,64)
			}
			if row.TxnType == "CREDIT"{
				cost ,_ =strconv.ParseFloat( row.CreditValue,64)
				//cost = row.CreditValue
			}
			item.Cost = cost
			statement.Total = statement.Total+ cost
			statement.ListItem  = append(statement.ListItem,item)
		}
	}
	return  statement
}

