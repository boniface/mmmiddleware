package cdn

import (
	"testing"
	"encoding/json"
	"fmt"

)

func TestCdnEntry_AddNewDoc(t *testing.T) {

	ref :=make(map[string]string)
	ref["orgcode"]="mm01"
	ref["department"]="sales"
	ref["docType"]="birth certificate"
	ref["inref"]="987kjuyt"
	ref["user"]="biangacila@gmail.com"
	ref["toten"]="ewwgsagsgsgggg9633"
	ref["module"]="qa"

	cdn :=CdnEntry{}
	cdn.Filename = "test01.txt"
	myref,_:= json.Marshal(ref)
	cdn.Ref = string(myref)
	cdn.Company = "nfumu"
	cdn.Campaign = "developer"
	cdn.ContentType = "text/plain"
	cdn.Department = "software"
	cdn.Email = "biangacila@gmail.com"
	cdn.Comment = "comment text"
	cdn.Reason = "reason comment"
	cdn.Inref = "zxcv9876541230"
	cdn.Status = "active"
	cdn.DocType = "text"
	/* save file into mongodb */
	cdn.SaveFileUrl()
	/* save file info into cassandra */
	cdn.AddNewDoc()
	docId :=cdn.MongoId
	fmt.Println("SAVED DOC ID > ",docId)
	/* Get doc from mongo */
	docInfo :=cdn.FndFileByIdGridfs(docId)
	str ,_:=json.Marshal(docInfo)
	fmt.Println("GET SAVED DOC BY ID > ",string(str))




}

func TestGetsavedfile(t *testing.T){
	cdn :=CdnEntry{}
	id :="57e68bb7b315f136a7b68eb3"
	/* Get doc from mongo */
	docInfo :=cdn.FndFileByIdGridfs(id)
	str ,_:=json.Marshal(docInfo)
	fmt.Println("GET SAVED DOC BY ID > ",string(str))
}
