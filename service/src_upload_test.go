package service

import (
	"testing"
	"fmt"
	"hashcode.zm/mmmiddleware/model"
	"github.com/pborman/uuid"


)

func TestCustomerUpload(t *testing.T) {
	sessionId :=uuid.New()
	referenceId :=uuid.New()
	referenceUploads := model.ReferenceUploads{}
	referenceUploads.Date ="2015-01-01"
	referenceUploads.Fullname ="biangacila merveilleux"
	referenceUploads.Login = "biangacila"
	referenceUploads.Username = "biangacila"
	referenceUploads.ReferenceId = referenceId
	referenceUploads.SessionId =sessionId
	referenceUploads.Url ="TBMAR10.csv" //"Jan 2015.csv" //
	//* todo get upload setting / for now let use our default setting for demo purpose
	upsetting := model.UploadSetting{}
	uploaSetting := upsetting.GetDefault()

	// now let get data from the url file
	filename:="../csv_uploaded/fresh/"+referenceUploads.Url
	filedata :=ReadCsvInto_CustomerUpload(filename)
	// now let get upload setting
	datefileinfo := model.DateInfoFile{}
	datefileinfo.Year, datefileinfo.Month, datefileinfo.Day, datefileinfo.Start, datefileinfo.End =FindDateYearMonthDay(filedata,uploaSetting)
	//now let build customer upload data array
	StartUploadService(filedata,uploaSetting,referenceUploads,datefileinfo)
	//fmt.Println("------->> DATA CUST >> ",custUpData)


}
func TestCheckIfFileInPool(t *testing.T) {
	listFile :=CheckIfFileInPool("test")
	if len(listFile) > 0{
		for _,filename :=range listFile{
			filedata :=ReadCsvInto_CustomerUpload(filename)
			fmt.Println("filedata => ",filedata)
			//FindDateYearMonthDay(filedata,)
		}
	}
}

func TestLoadDefaultFinanceCodePastel(t *testing.T){
	info :=LoadDefaultFinanceCodePastel()
	fmt.Println("Finance Statement info => ",len(info)," > ",info)
}

func TestFindFinanceStateInfo_Pastel(t *testing.T){
	mmcode:="5000/0000"
	fs,_ :=FindFinanceStateInfo_Pastel(mmcode)
	fmt.Println("catgory >>> "," > ",fs.Category)
	fmt.Println("SubCategory >>> "," > ",fs.SubCategory)
	fmt.Println("StartCode >>> "," > ",fs.StartCode)
	fmt.Println("EndCode >>> "," > ",fs.EndCode)
	fmt.Println("Type >>> "," > ",fs.Type)

}

//1473525269362