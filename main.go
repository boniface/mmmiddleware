package main

import (
	//"hashcode.zm/mmmiddleware/service"
	//"fmt"
	//"hashcode.zm/mmmiddleware/model"
	"hashcode.zm/mmmiddleware/cdn"
	"net/http"
	"log"
	ws "hashcode.zm/mmmiddleware/webservice"
)


func main() {

	/* set mongo db ipaddress */
	mongoIp :=cdn.GetMongoServerIp()
	cdn.MongoServerIp = mongoIp
	/* Run CDN server */
	go cdn.StartDataStore()
	/* RUN FOR EVER */
	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)
	http.HandleFunc("/uploaddoc", ws.UploadfileHandler)
	http.HandleFunc("/financiary/statement/income", ws.Ws_FinanciaryStatement_Income)
	http.HandleFunc("/financiary/statement/uploaded", ws.Ws_FinanciaryStatement_Uploaded)

	log.Println("Listening StaticWebsiteStore...")
	log.Fatal(http.ListenAndServe(":19004", nil))
}


