package service

import (
	"strings"
	"io"
	"bytes"
	"os"
	"log"
	"hashcode.zm/mmmiddleware/model"
	"fmt"
	"path/filepath"
	"strconv"
	"github.com/pborman/uuid"
	"hashcode.zm/mmmiddleware/dbconn"
	"github.com/fatih/structs"
	"hashcode.zm/communex-kiosk/generic"
)

type DateEF struct {
	Year  string
	Month string
	Day   string
}

func StartUploadService(filedata []string,uploadSetting model.UploadSetting,uploadInfo model.ReferenceUploads,dateinfo model.DateInfoFile)[]model.CustomerUpload {
	custUploadData :=[]model.CustomerUpload{}
	total :=0
	startX := uploadSetting.RowStart - 1
	//p :=fmt.Println
	for x,online :=range filedata{
		if x>=startX{
			row:=[]string{}
			if strings.Contains(online,";"){
				row =strings.Split(online,";")
			}
			if strings.Contains(online,","){
				row =strings.Split(online,",")
			}

			indexCode :=uploadSetting.CodeColom -1
			indexDesc :=uploadSetting.DescColom - 1
			indexDebt :=uploadSetting.DebitColom-1
			indexCredit :=uploadSetting.CreditColom-1

			mmCode := ""
			mmDesc :=""
			mmDebit :=""
			mmCredit :=""

			if contains(row, indexCode){
				mmCode = row[indexCode]
			}
			if contains(row, indexDebt){
				mmDebit = row[indexDebt]
			}
			if contains(row, indexCredit){
				mmCredit = row[indexCredit]
				mmCredit = strings.Trim(mmCredit,"\r")

			}
			if contains(row, indexDesc){
				mmDesc = row[indexDesc]
			}


			if mmCode !=""{
				txnType :="NA"
				if  mmCredit =="" && mmDebit ==""{
					txnType ="FEATURE"
				}else if mmCredit =="" && mmDebit !=""{
					txnType ="DEBIT"
				}else if mmCredit !="" && mmDebit ==""{
					txnType ="CREDIT"
				}


				if mmCredit ==""{
					mmCredit ="0"
				}
				if mmDebit ==""{
					mmDebit = "0"
				}

				mmDebit = strings.Trim(mmDebit," ")
				mmCredit = strings.Trim(mmCredit," ")

				debit ,_:=strconv.ParseFloat(mmDebit,64)
				credit ,_:=strconv.ParseFloat(mmCredit,64)

				custData :=builCustomerUploadData(online,uploadSetting,dateinfo,txnType,debit,credit,mmDesc,mmCode)
				custUploadData = append(custUploadData,custData)
				total++
			}

		}
	}
// select orgcode,year,month,AccountingCode,entrycategory,debitvalue,creditvalue,date  from custuploads where orgcode='MM01';
	return custUploadData


}
func builCustomerUploadData(rawdata string,uploadSetting model.UploadSetting,dateinfo model.DateInfoFile,txnType string,debitValue float64,creditValue float64, entryDesc string,mmcode string)model.CustomerUpload{
	//cdate,ctime :=generic.GetDateAndTimeString()
	d,t :=generic.GetDateAndTimeString()
	strDate := d+" "+t
	var err error
	fs,_ :=FindFinanceStateInfo_Pastel(mmcode)
	custUploadData :=model.CustomerUpload{}
	custUploadData.UpRef =""
	custUploadData.OrgCode =uploadSetting.OrgCode
	//custUploadData.Id =uuid.New()
	custUploadData.Refrence =""
	custUploadData.Date =strDate //   cdate+" "+ctime
	custUploadData.AccountingCode =mmcode
	custUploadData.Year,err =strconv.Atoi(dateinfo.Year)

	if err != nil {
		fmt.Println(" ERROR YEAR CANT BE EMPTY >>>> ",err," > ",dateinfo.Year)
		return custUploadData
	}
	custUploadData.Month,_ =strconv.Atoi(dateinfo.Month)
	custUploadData.Day,_ =strconv.Atoi(dateinfo.Day)
	custUploadData.AccountingSystem ="PASTEL"
	custUploadData.DebitValue = debitValue
	custUploadData.CreditValue =creditValue
	custUploadData.TxnType = txnType
	custUploadData.EntryDescription =entryDesc
	custUploadData.EntryCategory =fs.Category
	custUploadData.EntrySubCategory = fs.SubCategory
	custUploadData.CsvStringInput =""
	custUploadData.MappingCode =""
	custUploadData.DetermineTransctionType()

	if custUploadData.DebitValue ==0.00{
		custUploadData.DebitValue = 0
	}
	if custUploadData.CreditValue==0.00{
		custUploadData.CreditValue = 0
	}

	if custUploadData.EntryDescription !=""{
		if !strings.Contains(custUploadData.EntryDescription,"\r"){
			if !strings.Contains(custUploadData.EntryDescription,"\n"){
				mymap := structs.Map(custUploadData)
				cleanYear :=fmt.Sprintf("'%d'",custUploadData.Year)
				cleanMonth :=fmt.Sprintf("'%d'",custUploadData.Month)
				cleanDay :=fmt.Sprintf("'%d'",custUploadData.Day)
				rYear :=fmt.Sprintf("%d",custUploadData.Year)
				rMonth :=fmt.Sprintf("%d",custUploadData.Month)
				rday:=fmt.Sprintf("%d",custUploadData.Day)

				fZero :=fmt.Sprintf("'%d'",0)
				rZero :=fmt.Sprintf("%d",0)


				qry := dbconn.GetInsertQueryFromMapMM(mymap,"custuploads")
				qry = strings.Replace(qry,`0.00`,`0`,100)
				qry = strings.Replace(qry,`accounttingsystem`,`accountingsystem`,100)

				qry = strings.Replace(qry,cleanYear,rYear,100)
				qry = strings.Replace(qry,cleanMonth,rMonth,100)
				qry = strings.Replace(qry,cleanDay,rday,100)
				qry = strings.Replace(qry,fZero,rZero,100)

				dbconn.RunQueryCassCollection(qry)
			}
		}
	}
	return custUploadData
}
func CheckIfFileInPool(env string) []string {
	var finded []string
	dirname := "." + string(filepath.Separator) + "tmp-file-in/"

	if env == "test" {
		dirname = ".." + string(filepath.Separator) + "tmp-file-in/"
	}

	log.Println("******> ", dirname)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			fmt.Println(fi.Name(), fi.Size(), "bytes")
			targetDir := dirname + fi.Name()
			finded = append(finded, targetDir)
		}
	}
	log.Println("[data file list] >",len(finded))
	return finded
}
func  ReadCsvInto_CustomerUpload(filename string) []string {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filename) // Error handling elided for brevity.
	log.Println("Err filename ****> ", err)
	io.Copy(buf, f)           // Error handling elided for brevity.
	f.Close()
	s := string(buf.Bytes())
	s = strings.Replace(s,`'`,"",1000)
	/*
	let build the data now
	 */
	upList := []string{}
	lines := strings.Split(s, "\n")

	for _, oneline := range lines {
		//fmt.Println("===>", oneline)
		upList = append(upList, oneline)
	}
	return upList
}
func FindDateYearMonthDay(data []string,current model.UploadSetting) (string, string, string, string, string) {
	/*upsetting := model.UploadSetting{}
	current := upsetting.GetDefault()*/
	log.Println("current.DateColom----> ", current.DateColom)
	/*
	DATE START AND END ANALYSE
	 */
	dateStart :=DateEF{}
	dateEnd :=DateEF{}
	for x, row := range data {
		if x == (current.DateColom - 1) {
			dateStart,dateEnd =analizeDateFormation(current.DateFormat, row)
			//empty_space :=strings.Split(row,":")
			mystring := string(row)
			mystring = strings.Trim(mystring, " ")
		}
	}
	year := dateStart.Year
	month := dateStart.Month
	day := dateStart.Day
	start :=dateStart.Year+"-"+dateStart.Month+"-"+dateStart.Day
	end :=dateEnd.Year+"-"+dateEnd.Month+"-"+dateEnd.Day
	return year, month, day, start, end
}
func analizeDateFormation(format string, online string) (DateEF,DateEF) {
	arrStr := strings.Split(online, ":")
	//log.Println("arrStr)-->> ", len(arrStr))

	for xid, td := range arrStr {
		fmt.Println("[", xid, "] ", td)
	}
	targetOne := strings.Split(arrStr[1], " to ")
	dataStart :=extractDate(strings.Trim(targetOne[0]," "))
	dataEnd :=extractDate(strings.Trim(targetOne[1]," "))
	//find date one

return dataStart,dataEnd
}
func extractDate(datestring string) DateEF {
	//p :=fmt.Println
	arr := strings.Split(datestring, "/")
	mydate := DateEF{}
	if len(arr) >= 3 {
		mydate.Year = arr[2]
		mydate.Month = arr[1]
		mydate.Day = arr[0]

		mydate.Year = strings.Replace(mydate.Year,";","",20)
		mydate.Year = strings.Trim(mydate.Year," ")
		if len(mydate.Year)  ==2 || len(mydate.Year)  ==3 {
			mydate.Year = "20"+mydate.Year
		}
	}
	fmt.Println("=====>> Day: ",mydate.Day,"; Month: ",mydate.Month," ; year : ",mydate.Year)
	return mydate
}
var MapFinancialStatementInfo map[string]model.FinancialStatementInfo
func LoadDefaultFinanceCodePastel()map[string]model.FinancialStatementInfo{

	MapFinancialStatementInfo  = make(map[string]model.FinancialStatementInfo)
	fInfo :=make(map[string]model.FinancialStatementInfo)
	/* INCOME STATEMENT */
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"INCOME STATEMENT",Category:"SALES",SubCategory:"SALES", StartCode:1000,EndCode:1999}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"INCOME STATEMENT",Category:"COST OF SALES",SubCategory:"COST OF SALES", StartCode:2000,EndCode:2499}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"INCOME STATEMENT",Category:"OTHER INCOME",SubCategory:"OTHER INCOME", StartCode:2500,EndCode:2999}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"INCOME STATEMENT",Category:"EXPENSES",SubCategory:"EXPENSES", StartCode:3000,EndCode:4999}
	/* BALANCE SHEET */
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"BALANCE SHEET",Category:"EQUITY & LONG TERM LIABILITIES",SubCategory:"OWNERS EQUITY", StartCode:5100,EndCode:5499}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"BALANCE SHEET",Category:"EQUITY & LONG TERM LIABILITIES",SubCategory:"LONG TERM LIABILITIES", StartCode:5500,EndCode:5799}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"BALANCE SHEET",Category:"CURRENT ASSETS",SubCategory:"FIXED ASSETS", StartCode:6000,EndCode:6699}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"BALANCE SHEET",Category:"CURRENT ASSETS",SubCategory:"CURRENT ASSETS", StartCode:7000,EndCode:8999}
	fInfo[uuid.New()] =model.FinancialStatementInfo{Type :"BALANCE SHEET",Category:"CURRENT LIABILITIES",SubCategory:"CURRENT LIABILITIES", StartCode:9000,EndCode:9999}

	MapFinancialStatementInfo = fInfo

	return MapFinancialStatementInfo
}
func FindFinanceStateInfo_Pastel(mmcode string)(model.FinancialStatementInfo,bool){

	fsInfo :=model.FinancialStatementInfo{}
	boo :=false

	arrStr :=strings.Split(mmcode,"/")

	strCode:="0000"
	if len(arrStr)>0{
		strCode =arrStr[0]
	}


	/*re := regexp.MustCompile( "[^0-9]" )
	strCode = re.ReplaceAllString( strCode, "" )
	fmt.Printf( re.ReplaceAllString( "hello 12 34 five", "" ) )*/
	strCode= strings.Replace(strCode,"\n","",1)
	strCode= strings.Replace(strCode,"\r","",1)
	strCode = strings.Replace(strCode,";","",1)
	strCode = strings.Replace(strCode,",","",2)

	incode,err :=strconv.ParseInt(strCode,10,64)

	if err !=nil{
		log.Println("ERROR CONVERT TO INT64 CODE: ",err,strCode)
	}

	fsList := LoadDefaultFinanceCodePastel()

	for _,row :=range fsList{
		if checkIfNumBetween(incode,row.StartCode,row.EndCode){
			fsInfo = row
			boo =true
			break
		}
	}

	return fsInfo,boo
}
func checkIfNumBetween(num int64,startNum int64,endNum int64)bool{
	boo :=false
	if num >= startNum && num <= endNum{
		boo = true
	}
	return boo
}
func contains(s []string, e int) bool {
	mylen := len(s)
	if e >=0 && e <= mylen-1{
		return true
	}
	return false
}

/*
Mr b
 */