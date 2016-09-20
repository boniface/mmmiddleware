package cdn

import (
	"github.com/fatih/structs"
	"hashcode.zm/mmmiddleware/dbconn"
	"hashcode.zm/mmmiddleware/generic"
	"github.com/pborman/uuid"
)
// keyspace : marginmentor

type CdnEntry struct {
	Id          string
	Company     string
	Campaign    string
	Email       string
	Department  string
	Inref       string
	Extref      string
	DocType     string
	DocSize     string
	ContentType string
	Date        string
	Time        string
	MongoId     string
	Filename    string
	Status      string
	Ref string
	Reason  string
	Comment string
}

func (obj *CdnEntry)AddNewDoc(){
	obj.Date,obj.Time =generic.GetDateAndTimeString()
	obj.Id = uuid.New()
	mymap := structs.Map(obj)
	qryCdn := dbconn.GetInsertQueryFromMap(mymap,"kiosk_login")
	cdnCampaign :=dbconn.GetInsertQueryFromMap(mymap,"cdnCampaign")
	cdnEmail :=dbconn.GetInsertQueryFromMap(mymap,"cdnEmail")
	cdnDepartment:=dbconn.GetInsertQueryFromMap(mymap,"cdnDepartment")
	cdnDocType:=dbconn.GetInsertQueryFromMap(mymap,"cdnDocType")
	_,err:=dbconn.RunQueryCassCollectionFeedback(qryCdn)
	obj.CheckErr(err)
	_,err=dbconn.RunQueryCassCollectionFeedback(cdnCampaign)
	obj.CheckErr(err)
	_,err=dbconn.RunQueryCassCollectionFeedback(cdnEmail)
	obj.CheckErr(err)
	_,err=dbconn.RunQueryCassCollectionFeedback(cdnDepartment)
	obj.CheckErr(err)
	_,err=dbconn.RunQueryCassCollectionFeedback(cdnDocType)
	obj.CheckErr(err)
}
func(obj *CdnEntry)CheckErr(err error ){
	if err != nil {
		panic(err.Error())
	}
}
