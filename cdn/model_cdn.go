package cdn

import (
	"github.com/fatih/structs"
	"hashcode.zm/mmmiddleware/dbconn"
	"hashcode.zm/mmmiddleware/generic"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"io"
	"os"
	"fmt"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"path/filepath"
	"bytes"
	"log"
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
type FileMetaData struct {
	Id bson.RawD `sjon:"_id"`
	ChunkSize int64
	Filename string
	Length int64
	Md5 string
	UploadDate string
	ContentType string
	Metadata bson.Raw `bson:"metadata"`
}

type DataFile struct {
	Id string `json:"_id"`
	ChunkSize int64
	Filename string
	Length int64
	Md5 string
	UploadDate string
	ContentType string
	Metadata string `json:"metadata"`
}
type ChunkFile struct {
	Id  string `json:"_id"`
	Data []byte
	Files_id json.RawMessage  `json:"files_id"`

}
type ChunkFile2 struct {
	Id  string `json:"id"`
	Data []byte
	Files_id string

}

type DocFindInfo struct {
	Files DataFile
	Chunks ChunkFile
}

var  MongoServerIp string

func (obj *CdnEntry)AddNewDoc(){

	obj.SaveFileUrl()

	obj.Date,obj.Time =generic.GetDateAndTimeString()
	obj.Id = uuid.New()
	mymap := structs.Map(obj)
	qryCdn := dbconn.GetInsertQueryFromMap(mymap,"cdn")
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

func(obj *CdnEntry)SaveFileUrl() {
	filename :=obj.Filename
	dirUrl :="./tmp-file-in/"+filename
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic("Saving file Error cdn > "+err.Error() +" > "+filename)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(obj.Company)
	file, err := db.GridFS("fs").Create(filename)
	check(err)
	messages, err := os.Open(dirUrl)
	check(err)
	defer messages.Close()
	_, err = io.Copy(file, messages)
	check(err)
	/* Let construct the meta data for cassandra info */
	id :=fmt.Sprintf("%v",file.Id())
	id = strings.Replace(id,`ObjectIdHex("`,"",1)
	id = strings.Replace(id,`")`,"",1)
	err = file.Close()
	check(err)
	obj.MongoId = id
}

func(obj *CdnEntry) FndFileByIdGridfs(id string)DocFindInfo{
	if MongoServerIp ==""{
		MongoServerIp = GetMongoServerIp()
	}
	docinfo :=DocFindInfo{}
	session, err := mgo.Dial(MongoServerIp)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	/* let find the thunk file */
	chunkFile :=ChunkFile{}
	db2 := session.DB(obj.Company).C("fs.chunks")
	var myinterface  interface{}
	db2.Find( bson.M{"files_id":bson.ObjectIdHex(id)}).One(&myinterface)
	strout ,err:=json.Marshal(myinterface)
	json.Unmarshal(strout,&chunkFile)
	docinfo.Chunks = chunkFile
	/* end of that file */
	db3 := session.DB("test").C("fs.files")
	var dataF interface{}
	db3.Find( bson.M{"_id":bson.ObjectIdHex(id)}).One(&dataF)
	strout2 ,err:=json.Marshal(dataF)
	json.Unmarshal(strout2,&dataF)
	str,_:=json.Marshal(dataF)
	dataFile :=DataFile{}
	json.Unmarshal(str,&dataFile)
	docinfo.Files = dataFile
	return docinfo
}

func GetMongoServerIp()string{
	dirname := "." + string(filepath.Separator) + "mongo.conf"
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		// dirname does not exist
		dirname = ".." + string(filepath.Separator) + "mongo.conf"
	}
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(dirname)
	if err != nil {
		log.Println("Err filename ****> ", err)
		return "0.0.0.0"
	}

	io.Copy(buf, f)
	f.Close()

	s := string(buf.Bytes())

	/*
	let build the data now
	 */
	upList := []string{}
	lines := strings.Split(s, "\n")
	for _, oneline := range lines {
		upList = append(upList, oneline)
	}
	log.Println("CASSANDRA COMM REQ +++++> ", upList)
	d.Close()

	Host1:="0.0.0.0"

	if len(upList) > 0 {
		Host1 = upList[0]
	}

	return Host1

}
