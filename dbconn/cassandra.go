package dbconn

import (
	"github.com/gocql/gocql"
	"fmt"
	"bytes"
	"os"
	"log"
	"io"
	"strings"
	"reflect"
	"encoding/json"
)

var Cass_keyspace string = "marginmentor"
var CassSession *gocql.Session

var Host1 = "127.0.0.1"


func CreateKeyspace(){
	qry:=`CREATE KEYSPACE `+Cass_keyspace+` WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};`
	cluster := gocql.NewCluster(Host1)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	inner_session, errCass := cluster.CreateSession()
	defer inner_session.Close()
	if errCass != nil {
		fmt.Println("CreateKeyspace ERROR > ", errCass)
	}
	if err := inner_session.Query(qry).Exec(); err != nil {
		fmt.Println("CreateKeyspace Error --->>> ", err)
	}
	inner_session.Close()

	/* let create table baased on the cql file */
	cluster.Keyspace = Cass_keyspace
	inner_session, errCass = cluster.CreateSession()
	session := inner_session
	qrys := GetDbSqlPing()
	fmt.Println("[MM SCHEMA QRYS] > ", qrys)
	arrQry := strings.Split(qrys, ";")
	for _, qry := range arrQry {
		if err := session.Query(qry).Exec(); err != nil {
			//log.Fatal(err)
			fmt.Println("--->>> ", err)
		}
	}

	session.Close()
}
func GetDbSqlPing() string {
	filename := ""
	if _, err := os.Stat("./mm_schema.sql"); err == nil {
		// path/to/whatever exists
		filename = "./mm_schema.sql"
	}else{
		filename = "../mm_schema.sql"
	}

	filenames := []string{filename}
	buf := bytes.NewBuffer(nil)
	for _, filename := range filenames {
		f, err := os.Open(filename) // Error handling elided for brevity.
		log.Println("Err filename ****> ", err)
		io.Copy(buf, f)           // Error handling elided for brevity.
		f.Close()
	}
	s := string(buf.Bytes())
	return s
}
func InsertRowTable(obj interface{},table string){
	qry :=QryInsertIntoStructTable(obj,table)
	cluster := gocql.NewCluster(Host1)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.Keyspace = Cass_keyspace
	inner_session, errCass := cluster.CreateSession()
	defer inner_session.Close()
	if errCass != nil {
		fmt.Println("CreateKeyspace ERROR > ", errCass)
	}
	if err := inner_session.Query(qry).Exec(); err != nil {
		fmt.Println("CreateKeyspace Error --->>> ", err)
	}
	inner_session.Close()
}
func RunQueryCassCollection(qry string)string{
	cluster := gocql.NewCluster(Host1)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.Keyspace = Cass_keyspace
	session, errCass := cluster.CreateSession()
	defer session.Close()
	if errCass != nil {
		fmt.Println("ERROR  RunQueryCassCollection> ", errCass)
	}
	if err := session.Query(qry).Exec(); err != nil {
		fmt.Println("CreateKeyspace Error --->>> ", err)
	}
	iter :=session.Query(qry).Iter()
	myrow ,_:=iter.SliceMap()
	str,_ :=json.Marshal(myrow)
	session.Close()
	return string(str);
}
func QryInsertIntoStructTable(mystruct interface{},table string) string{
	v := reflect.ValueOf(mystruct)
	qry :=""
	col :="("
	valcol :="("
	nunCol :=0
	typ :=  v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i).Name
		fieldCol :=fmt.Sprintf("%v",f)
		fieldValue :=fmt.Sprintf("%v",v.Field(i).Interface())
		if nunCol ==0{
			col =col+""+fieldCol+""
			valcol = valcol+"'"+fieldValue+"' "
			nunCol++
		}else{
			col =col+","+fieldCol+""
			valcol = valcol+",'"+fieldValue+"' "
			nunCol++
		}
	}
	col =col+")"
	valcol = valcol+") "
	qry = "insert into "+table+""+col+" values"+valcol
	return qry
}
