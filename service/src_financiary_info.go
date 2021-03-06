package service

import (
	"hashcode.zm/mmmiddleware/model"
	"hashcode.zm/mmmiddleware/dbconn"
	"encoding/json"
	"fmt"
)

type UploadedInfo struct {
	Year int
	Month int
	Day int
	Username string
	Filename string
	DateTime string
	OrgCode string
}
type FinanceInfo struct {
	ReqType string
	Year int
	Month int
	OrgCode string
	UploadedInfo []UploadedInfo
}

func(uInfo *FinanceInfo)LoadOneYeaUploadedInfo(){
	ls := []model.CustomerUpload{}
	qry :=fmt.Sprintf("select * from custUploads where orgcode ='%s'  order by year ASC, month ASC ",uInfo.OrgCode)
	fmt.Println(qry)
	strRs :=dbconn.RunQueryCassCollection(qry)
	json.Unmarshal([]byte(strRs),&ls)
	mymap :=make(map[string]model.CustomerUpload)
	for _,row:=range ls{
		key :=fmt.Sprintf("%d#%d#%d",row.OrgCode,row.Year,row.Month)
		mymap[key]=row
	}

	for _,row :=range mymap{
		info :=UploadedInfo{}
		info.Year = row.Year
		info.Month =row.Month
		info.Day = row.Day
		info.Username = row.Date
		uInfo.UploadedInfo = append(uInfo.UploadedInfo,info)
	}

}

// select * from custUploads where orgcode ='MM01' and year=2001

// select orgcode,year,month,accountingcode,entrycategory,entrydescription,debitvalue,creditvalue,date from custuploads where orgcode='MM01' and year=2011 ;
