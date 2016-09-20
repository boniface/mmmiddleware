package cdn

import (

	"net/http"
	"os"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
	"io"
	"strings"
	"log"

)

// Acceess handler
func Access(res http.ResponseWriter, req *http.Request) {
	// check session or something like this
}

// Set prefix for current collection name, of course, if need it
// It useful when you store another data in one database
func AutoPrefix(params martini.Params) {
	// 'coll' - parameter name, which is used
	params["coll"] = "cdn." + params["coll"]
	// Ok. Now, cdn will work with this prefix for all collections
}


type Person struct {
	Name  string
	Phone string
}



func check(err error) {
	if err != nil {
		panic(err.Error())
	}
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

type FileMetaData struct {
	Id bson.RawD `bson:"_id"`
	ChunkSize int64
	Filename string
	Length int64
	Md5 string
	UploadDate string
	ContentType string
	Metadata bson.Raw `bson:"metadata"`
}

type DocId struct {
	Ref     string
	Company string
}

func SaveFileByte() string {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	file, err := db.GridFS("fs").Create("myfile2.txt")
	check(err)
	n, err := file.Write([]byte("Hello world 222!"))
	check(err)
	id := file.Id()
	err = file.Close()
	check(err)
	fmt.Printf("%d bytes written\n -=>%4", n, id)
	var result bson.M
	err = db.C("fs").Find(nil).One(&result)

	type II struct {
		Id string
	}

	ii := II{}
	str, _ := json.Marshal(result)
	json.Unmarshal(str, &ii)
	return ii.Id
}

func FindAllDoc() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	iter := db.GridFS("chunks").Find(nil).Iter()
	project :=bson.M{}
	for iter.Next(project) {
		str, err := json.Marshal(project)
		fmt.Println("--:) ", err, string(str))
	}
}

func SaveFileUrl(dirUrl string, filename string,ref map[string]string) string {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")

	file, err := db.GridFS("fs").Create(filename)
	check(err)
	messages, err := os.Open(dirUrl)
	check(err)
	defer messages.Close()
	var resultSave int64
	resultSave, err = io.Copy(file, messages)
	check(err)

	/* Let construct the meta data for cassandra info */

	id :=fmt.Sprintf("%v",file.Id())
	id = strings.Replace(id,`ObjectIdHex("`,"",1)
	id = strings.Replace(id,`")`,"",1)

	contentType :=file.ContentType()
	//metadata :=GetMetadata(id)
	md5 :=file.MD5()

	fmt.Println("==saved id ===> ",id," > ",contentType," > "," > ",md5)

	err = file.Close()
	check(err)
	str := fmt.Sprintf("Done saving file err: %v, writen: %v", err, resultSave)
	return str
}

func FindChunksFile(file_id string)ChunkFile2{
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test").C("fs.chunks")
	var myinterface  []interface{}
	db.Find(nil).All(&myinterface)
	str,err :=json.Marshal(myinterface)
	mychunks :=[]ChunkFile{}
	mychunk :=ChunkFile2{}
	err = json.Unmarshal(str,&mychunks)
	listFiles :=[]ChunkFile2{}
	json.Unmarshal(str,&listFiles)
	for _,row :=range listFiles{
		if row.Files_id == file_id{
			fmt.Println("--:) "," >> "," > ",row.Files_id," > ",row.Id)
			mychunk = row
		}
	}
	return mychunk
}

func GetMetadata(strId string)FileMetaData{
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test").C("fs.files")
	var myinterface  []FileMetaData
	db.Find(nil).All(&myinterface)
	str,err :=json.Marshal(myinterface)
	fmt.Println("FINDED FILES > ",len(myinterface),myinterface)
	if err !=nil{
		log.Println("@@@@ 1: ",err,str)
	}
	metadata :=FileMetaData{}

	/*check(err)
	ls :=[]FileMetaData{}
	metadata :=FileMetaData{}

	err = json.Unmarshal(str,&ls)
	if err !=nil{
		log.Println("@@@@ 2: ",err,string(str))
	}*/
	//check(err)
	for _,row:=range myinterface{
		id :=fmt.Sprintf("%v",row.Id)
		id = strings.Replace(id,`ObjectIdHex("`,"",1)
		id = strings.Replace(id,`")`,"",1)
		if id == strId{
			metadata = row
			fmt.Println("--:) ", err, row)
		}
	}
	fmt.Println("<<--:) ", err, string(str))
	return metadata
}

/*func Tmain() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
func Rummain() {
	m := martini.Classic()

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("cdn")
	m.Map(db)

	logger := log.New(os.Stdout, "\x1B[36m[cdn] >>\x1B[39m ", 0)
	m.Map(logger)

	*//*m.Group("/uploads",cdn.
		cdn.Cdn(cdn.Config{
			// Maximum width or height with pixels to crop or resize
			// Useful to high performance
			MaxSize:  1000,
			// Show statictics and the listing of files
			ShowInfo: true,
			// If true it send URL without collection name, like this:
			// {"field":"/5364d634952b829316000001/books.jpg", "error": null}
			TailOnly: true,
		}),
		// Access logic here
		Access,
		// On the fly prefix for collection
		AutoPrefix,
	)*//*

	logger.Println("Server started at :3000")
	m.Run()
}*/
