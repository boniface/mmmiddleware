package webservice

import (
	"net/http"
	"hashcode.zm/mmmiddleware/service"
	"encoding/json"
	"strconv"
)

func Ws_FinanciaryStatement_Uploaded(w http.ResponseWriter, r *http.Request){
	orgcode :=r.Form["orgcode"][0]
	year :=r.Form["year"][0]
	montth :=r.Form["month"][0]
	ReqType :=r.Form["reqtype"][0]

	info :=service.FinanceInfo{}
	info.Year ,_=strconv.Atoi( year)
	info.Month ,_= strconv.Atoi(montth)
	info.OrgCode =orgcode
	info.ReqType =ReqType
	info.LoadOneYeaUploadedInfo()
	rows :=info.UploadedInfo

	mybyte ,_:=json.Marshal(rows)
	w.Write(mybyte)

}
func Ws_FinanciaryStatement_Income(w http.ResponseWriter, r *http.Request){

}
