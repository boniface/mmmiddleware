package cdn

import (
	"log"
	"net/http"
	"os"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
	"github.com/pborman/uuid"
	"time"
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

	/*m.Group("/uploads",cdn.
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
	)*/

	logger.Println("Server started at :3000")
	m.Run()
}

type Person struct {
	Name  string
	Phone string
}

func Tmain() {
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

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type M struct {
	Id          interface{} "_id"
	ChunkSize   int         "chunkSize"
	UploadDate  time.Time   "uploadDate"
	Length      int64       ",minsize"
	MD5         string
	Filename    string    ",omitempty"
	ContentType string    "contentType,omitempty"
	Metadata    *bson.Raw ",omitempty"
	/*Id        string
	Length     int
	ChunkSize  float64
	UploadDate time.Time
	Md5        string
	Filename   string
	Data []byte
	Files_id string
	ContentType string
	Metadata string
	Any string*/
}

type DocId struct {
	Ref     string
	Company string
}

func SaveData() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("test")

	file, err := db.GridFS("fs").Create("myfile.txt")

	check(err)

	n, err := file.Write([]byte("Hello world!"))
	check(err)

	var info interface{}
	err = file.GetMeta(&info)
	myid := uuid.New()
	var docId DocId
	docId.Ref = uuid.New()
	docId.Company = "africa direct"
	file.SetId(docId)
	file.SetName("myfile2.txt")
	file.SetContentType("text/plain")
	//file.SetMeta(M{Any: "BiaLuv00123"})

	t := time.Date(2014, 1, 1, 1, 1, 1, 0, time.Local)
	file.SetUploadDate(t)

	err = file.Close()
	check(err)
	fmt.Printf("%d bytes written\n", n)

	var result M
	err = db.C("fs.files").Find(nil).One(&result)

	coFiles, errOpen := db.GridFS("fs").Open("myfile22.txt")

	str, _ := json.Marshal(result)
	str2, _ := json.Marshal(coFiles)

	fmt.Println(myid, " >> ", errOpen, string(str2))
	fmt.Println(string(str))
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
	fmt.Printf("%d bytes written\n -=>%4", n,id)
	var result M
	err = db.C("fs.files").Find(nil).One(&result)

	type II struct {
		Id string
	}

	ii := II{}
	str, _ := json.Marshal(result)
	json.Unmarshal(str, &ii)
	return ii.Id
}

func OpenSavedFile(fileId string) {
	// 57dec1b3b315f13eb095d0b8
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")

	//myid := interface{"_id":fileId,}
	file, err := db.GridFS("fs").OpenId(fileId)  //Create("myfile2.txt")
	log.Println("/***error**/>> ", err)
	if err != nil {
		os.Exit(1)
	}

	b := make([]byte, 30)
	n, err := file.Read(b)

	fmt.Println("/*****/>> ", n, string(b))

	err = file.Close()

}