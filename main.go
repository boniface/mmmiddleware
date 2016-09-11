package main

import (
	"hashcode.zm/mmmiddleware/service"
	"fmt"
	"hashcode.zm/mmmiddleware/model"
)


func main() {

	// demo only
	//todo please switch to realy upload setting late
	upsetting := model.UploadSetting{}
	current := upsetting.GetDefault()

	p :=fmt.Println
	listFile := service.CheckIfFileInPool("live")
	strdis :=fmt.Sprintf(">>>>filedata => %d ",len(listFile))
	p(strdis)
	if len(listFile) > 0{
		p(">>PASS 01")
		for _,filename :=range listFile{
			datefileinfo :=model.DateInfoFile{}
			p(">>PASS 02")
			filedata :=service.ReadCsvInto_CustomerUpload(filename)
			p(">>PASS 03 ROWS:  ",len(filedata))
			//fmt.Sprint  f("filedata => %v ",filedata)
			datefileinfo.Year, datefileinfo.Month, datefileinfo.Day, datefileinfo.Start, datefileinfo.End =service.FindDateYearMonthDay(filedata,current)
			p(">>PASS 04 > ",datefileinfo.Year, datefileinfo.Month, datefileinfo.Day, datefileinfo.Start)

		}
	}


}
