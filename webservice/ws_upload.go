package webservice

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"time"
	"crypto/md5"
	"strconv"
	"html/template"
	"log"
	"github.com/pborman/uuid"
	"hashcode.zm/mmmiddleware/model"
	"encoding/json"
	"hashcode.zm/mmmiddleware/service"
	"github.com/fatih/structs"
	"hashcode.zm/mmmiddleware/dbconn"
	"html"
	"hashcode.zm/mmmiddleware/cdn"
)

// Acceess handler
func UploadfileHandler(w http.ResponseWriter, r *http.Request) {
	uploadedUrl :="./tmp-file-in/"
	// accept orgin
	OrignAllowed(w)


	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		//contentfile :=r.MultipartForm()
		file, handler, err := r.FormFile("file")

		if err != nil {
			fmt.Println("file, handler, err > ", err, " > ", file, handler, r.MultipartForm)
			//return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Println("@@@--> ", uploadedUrl + handler.Filename)
		f, err := os.OpenFile(uploadedUrl + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			//return
		}
		defer f.Close()
		io.Copy(f, file)

		go ProcessUploadCustomerLive(w,r,handler.Filename)
	}





}
func ProcessUploadCustomerLive(w http.ResponseWriter,r *http.Request,filename string){
	dir :="./tmp-file-in/"
	referenceId :=uuid.New()
	referenceUploads := model.ReferenceUploads{}
	referenceUploads.OrgCode= r.Form["orgcode"][0]
	referenceUploads.Date =r.Form["date"][0] //"2015-01-01"
	referenceUploads.Fullname =r.Form["fullname"][0] //"biangacila merveilleux"
	referenceUploads.Login = r.Form["login"][0] //"biangacila"
	referenceUploads.Username =r.Form["username"][0]// "biangacila"
	referenceUploads.ReferenceId = referenceId
	referenceUploads.SessionId = r.Form["sessionid"][0]
	referenceUploads.Url =filename//"Jan 2015.csv" //
	str ,err := json.Marshal(referenceUploads)
	log.Println("+++++++ Posted Data To uplaod > ",err,string(str))

	// now let save the file into our cdn mongo db
	mymap := structs.Map(referenceUploads)
	cdnId :=SaveFileIntoCdn(mymap , referenceUploads)
	log.Println("######>> Finish SaveFileIntoCdn > ",cdnId)

	//* todo get upload setting / for now let use our default setting for demo purpose
	upsetting := model.UploadSetting{}
	uploaSetting := upsetting.GetDefault()
	/* add upsetting to reference upload */
	referenceUploads.UploadSettingsId = uploaSetting.Id
	// now let get data from the url file
	filename2:=dir+referenceUploads.Url
	filedata := service.ReadCsvInto_CustomerUpload(filename2)
	log.Println("######>> Finish ReadCsvInto_CustomerUpload > ")

	// now let get upload setting
	datefileinfo := model.DateInfoFile{}
	datefileinfo.Year, datefileinfo.Month, datefileinfo.Day, datefileinfo.Start, datefileinfo.End =service.FindDateYearMonthDay(filedata,uploaSetting)
	log.Println("[DATE INFO ] > ",datefileinfo.Year, datefileinfo.Month, datefileinfo.Day)
	//now let build customer upload data array
	service.StartUploadService(filedata,uploaSetting,referenceUploads,datefileinfo)

	//let now save upload reference / referenceuploads
	referenceUploads.Url = cdnId
	mymap = structs.Map(referenceUploads)
	qry := dbconn.GetInsertQueryFromMap(mymap,"referenceuploads")
	dbconn.RunQueryCassCollection(qry)
	log.Println("######>> Finish upload reference  > ")
	//defer CloseMyCon(w,r)
	fmt.Fprintf(w, "POST, %q", html.EscapeString(r.URL.Path))
}
func CloseMyCon(w http.ResponseWriter,r *http.Request){
	str2 := "done"
	mybyte :=[]byte(str2)
	w.Write(mybyte)
	err :=r.Body.Close()
	fmt.Println("we close the connection for this client > ",err)
}
func SaveFileIntoCdn(mymap map[string]interface{},upRef  model.ReferenceUploads)string{
	filepath :="./tmp-file-in/"+upRef.Url
	file, err := os.Open( filepath )
	if err != nil {
		log.Fatal("os.Open log.Fatal > ",err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)/* server static content server */
	}
	fmt.Println( fi.Size() )


	cdn :=cdn.CdnEntry{}
	cdn.Filename = upRef.Url //"test01.txt"
	myref,_:= json.Marshal(mymap)
	cdn.Ref = string(myref)
	cdn.Company = upRef.OrgCode //"nfumu"
	cdn.Campaign = "MM CUSTOMER UPLOAD"
	cdn.ContentType ="text/csv"
	cdn.Department = "software"
	cdn.Email =upRef.Username // "biangacila@gmail.com"
	cdn.Comment = "NO COMMENT"
	cdn.Reason = "NO REASON"
	cdn.Inref =upRef.ReferenceId//"zxcv9876541230"
	cdn.Status = "active"
	cdn.DocType = "csv"
	/* save file into mongodb */
	cdn.SaveFileUrl()
	/* save file info into cassandra */
	cdn.AddNewDoc()
	docId :=cdn.MongoId
	fmt.Println("SAVED DOC ID > ",docId)
	return string(docId)
}
func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func authorizeheader() string {
	str := "Origin,X-Requested-With,Accept, X-Token, x-token,Content-Type,X-Custom-Header,ics_email,ics_fullname,ics_phone,ics_role,ics_token,ics_username"
	return str
}
func OrignAllowed(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", authorizeheader())
	w.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")
}

