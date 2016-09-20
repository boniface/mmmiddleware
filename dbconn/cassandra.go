package dbconn

import (
	"github.com/gocql/gocql"
	"fmt"
	"strings"
	"reflect"
	"path/filepath"
	"os"
	"bytes"
	"log"
	"io"
	"encoding/json"
)

var Cass_keyspace string = "marginmentor"
var CassSession *gocql.Session

var Host1 = "172.17.0.2"
var Host2 = "172.17.0.5" //"41.193.232.157"
var Host3 = "172.17.0.4"
var boo bool = false


func init(){
	CreateCassandraSession2()
}
func CreateKeyspace() {
	qry := `CREATE KEYSPACE ` + Cass_keyspace + ` WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};`
	cluster := gocql.NewCluster(Host1)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.ProtoVersion = 3
	inner_session, errCass := cluster.CreateSession()
	defer inner_session.Close()
	if errCass != nil {
		fmt.Println("CreateKeyspace ERROR > ", errCass)
	}
	if err := inner_session.Query(qry).Exec(); err != nil {
		fmt.Println("CreateKeyspace Error --->>> ", err)
	}
}

func QryUpdateStructTable(mystruct interface{}, table string, avoided []string, conditions map[string]string) string {
	v := reflect.ValueOf(mystruct)
	qry := "update " + table + " set "
	nunCol := 0
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i).Name
		fieldCol := fmt.Sprintf("%v", f)
		fieldValue := fmt.Sprintf("%v", v.Field(i).Interface())
		if !IsInArray(avoided,fieldCol ) {
			if nunCol == 0 {
				qry = qry + " " + fieldCol + " = '" + fieldValue + "' "
				nunCol++
			} else {
				qry = qry + ", " + fieldCol + " = '" + fieldValue + "' "
				nunCol++
			}
		}

	}

	x := 0
	for condkey, condvalue := range conditions {
		if x == 0 {
			qry = qry + " where " + condkey + " = '" + condvalue + "' "
			x++
		} else {
			qry = qry + " AND " + condkey + " = '" + condvalue + "' "
			x++
		}
	}
	return qry
}

func QryInsertIntoStructTable(mystruct interface{}, table string) string {
	v := reflect.ValueOf(mystruct)
	qry := ""
	col := "("
	valcol := "("
	nunCol := 0
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i).Name
		fieldCol := fmt.Sprintf("%v", f)
		fieldValue := fmt.Sprintf("%v", v.Field(i).Interface())
		if nunCol == 0 {
			col = col + "" + fieldCol + ""
			valcol = valcol + "'" + fieldValue + "' "
			nunCol++
		} else {
			col = col + "," + fieldCol + ""
			valcol = valcol + ",'" + fieldValue + "' "
			nunCol++
		}
	}
	col = col + ")"
	valcol = valcol + ") "
	qry = "insert into " + table + "" + col + " values" + valcol
	return qry
}
func CreateCassandraSession2() {
	SetDbServerFilename()
	slist := GetServerCassandra("live")
	cluster := gocql.NewCluster(slist[0], slist[1], slist[2])
	cluster.Keyspace = Cass_keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.ProtoVersion = 3
	inner_session, errCass := cluster.CreateSession()
	//defer inner_session.Close()
	if errCass != nil {
		fmt.Println("errCass > ", errCass)
		//return errCass
	}
	//fmt.Println("CONNECTION TO CASSANDRA 100%")
	CassSession = inner_session

}
func CreatTabs() {
	session := CassSession
	for _, qry := range DbTablesQuery() {
		if err := session.Query(qry).Exec(); err != nil {
			//log.Fatal(err)
			fmt.Println("--->>> ", err)
		}
	}
}
func DbTablesQuery() []string {

	mystring :=
	`CREATE TABLE users(
  id text,
  name text,
  surname text,
  phone text,
  email text,
  username text ,
  address text,
  idnumber text,
  type text,
  role text,
  level text,
  status text,
  created text ,
  modified text,
  PRIMARY KEY (email)
);

CREATE TABLE users_login(
  id text,
  email text,CREATE TABLE tickets_comment(
  id text,
  tid text,
  kiosk text,
  username text,
  typeticket text,
  name text,
  date text,
  time text,
  status text ,
  comment text,
  title text,
  reason text,
  level text,
  PRIMARY KEY (kiosk,date,typeticket,tid,id)
);

CREATE TABLE users_status(
  id text,
  old_status text,
  new_status text,
  change_by text,
  email text,
  username text ,
  created text ,
  PRIMARY KEY (email,id)
);
CREATE TABLE ambassador_kiosk(
  id text,
  email text ,
  username text,
  kiosk text ,
  Date text,
  Time text,
  assign_by text ,
  status text ,
  PRIMARY key(email,id)
);
CREATE TABLE ambassador_kiosk_history(
  id text,
  email text ,
  username text,
  kiosk text ,
  date text,
  time text,
  assign_by text ,
  status text ,
  PRIMARY key(email,id)
);
CREATE TABLE tickets(
  tid text ,
  typeticket text,
  kiosk text,
  ambassador text,
  username text,
  date text,
  time text,
  name text,
  standard text ,
  catogory text ,
  level text,
  input text,
  expected_complete_date text,
  expected_complete_time text,
  status text ,
  title text,
  reason text,
  comment text,
  over_due text ,
  thread text,
  priority text ,
  notification_thread text,
  expired text ,
  expired_date text,
  expired_time text,
  created text,
  PRIMARY key(tid)
);
CREATE TABLE tickets_ambassador(
  tid text,
  kiosk text,
  ambassador text,
  typeticket text,
  name text,
  date text,
  time text,
  PRIMARY KEY (ambassador,date,typeticket,tid)
);
CREATE TABLE tickets_kiosk(
  tid text,
  kiosk text,
  ambassador text,
  typeticket text,
  name text,
  date text,
  PRIMARY KEY (kiosk,date,typeticket,tid)
);
CREATE TABLE tickets_status(
  id text,
  tid text,
  kiosk text,
  username text,
  typeticket text,
  name text,
  date text,
  time text,
  status text ,
  modified text,
  PRIMARY KEY (tid,status,date,id)
);
CREATE TABLE tickets_comment(
  id text,
  tid text,
  kiosk text,
  username text,
  typeticket text,
  name text,
  date text,
  time text,
  status text ,
  comment text,
  title text,
  reason text,
  level text,
  PRIMARY KEY (tid,date,status,id)
);
CREATE TABLE tickets_escalation(
  id text,
  tid text,
  kiosk text,
  username text,
  typeticket text,
  name text,
  date text,
  time text,
  status text ,
  level text ,
  PRIMARY KEY (level,date,status,kiosk,tid,id)
);
`
	str := strings.Split(mystring, ";")

	return str
}

var dbserverFilename string

var UpList []string

func RunQueryCassCollectionFeedback(qry string) (string,error) {
	slist := GetServerCassandra("live")
	if len(slist) ==0{
		SetDbServerFilename()
		slist = GetServerCassandra("live")
	}
	err := CassSession.Query(qry).Exec();
	fmt.Println("RunQueryCassCollection Error --->>> ", err, " > ", qry)
	/*if  err != nil {
		fmt.Println("RunQueryCassCollection Error --->>> ", err, " > ", qry)
	}*/
	iter := CassSession.Query(qry).Iter()
	myrow, _ := iter.SliceMap()
	str, _ := json.Marshal(myrow)


	return string(str),nil;
}
func RunQueryCassCollection(qry string) string {

	slist := GetServerCassandra("live")
	if len(slist) ==0{
		SetDbServerFilename()
		slist = GetServerCassandra("live")
	}

	if err := CassSession.Query(qry).Exec(); err != nil {
		fmt.Println("RunQueryCassCollection Error --->>> ", err, " > ", qry)
	}
	iter := CassSession.Query(qry).Iter()
	myrow, _ := iter.SliceMap()
	str, _ := json.Marshal(myrow)

	return string(str);
}
func SetDbServerFilename() []string {
	dirname := "." + string(filepath.Separator) + "dbserver.init"
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		// dirname does not exist
		dirname = ".." + string(filepath.Separator) + "dbserver.init"
	}
	dbserverFilename = dirname

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
		return []string{}
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
	UpList = upList

	log.Println("CASSANDRA COMM REQ +++++> ", UpList)

	d.Close()

	if len(upList) > 0 {
		Host1 = upList[0]
	}

	return upList
}
func GetServerCassandra(env string) []string {
	return UpList
}
func InitCreateDatabaseTablesSCHEMA() {

	sl := GetServerCassandra("live")
	server := "127.0.0.1"
	if len(sl) > 0 {
		server = sl[0]
	}
	cluster := gocql.NewCluster(server)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.Keyspace = "ckiosk"
	cluster.ProtoVersion = 3
	inner_session, errCass := cluster.CreateSession()
	//defer inner_session.Close()
	if errCass != nil {
		fmt.Println("errCass > ", errCass)
		//return errCass
	}
	//fmt.Println("CONNECTION TO CASSANDRA 100%")
	CassSession = inner_session

	session := CassSession
	cluster.ProtoVersion = 3

	qry := `
create table onlineckioskreading(
  Id text,
  Hid text,
	Ambassador text,
	KioskName text,
	Ipaddr text,
	MemTotal text,
	MemUsed text,
	MemAvailable text,
	CpuTotal text,
	CpuUsed text,
	CpuAvailable text,
	LatencyServer text,
	LatencyApi text,
	TicketOpen text,
	TickerCLose text,
	TicketEscalate text,
	TicketStatusRed text,
	TicketStatusYallow text,
	TicketStatusBlue text,
	TicketStatusGreen text,
	TicketStatusOrange text,
	Date text,
	Time text,
	PRIMARY KEY (hid,date,id)
);`

	if err := session.Query(qry).Exec(); err != nil {
		//log.Fatal(err)
		fmt.Println("--->>> ", err)
	}

}

func GetInsertQueryFromMap(mymap map[string]interface{}, table string) string {
	x := 0
	var col, values string
	for key, val := range mymap {
		strval := ""

		if _, ok := val.(int64); ok {
			d := fmt.Sprintf("%d", val.(int64))
			strval = d
		}
		if _, ok := val.(float64); ok {
			d := fmt.Sprintf("%.2f", val.(float64))
			strval = d
		}
		if _, ok := val.(string); ok {
			strval = string(val.(string))
		}

		if x == 0 {
			if strval != "" {
				col = col + "" + key + " "
				values = values + "'" + strval + "' "
				x++
			}
		} else {
			if strval != "" {
				col = col + "," + key + " "
				values = values + ",'" + strval + "' "
				x++
			}

		}

	}
	var str string = "insert into " + table + "(" + col + ") values(" + values + ") ";
	return str
}
func GetUpdateQueryFromMap_ConditionEqual(table string, mymapdata map[string]interface{}, mymapcond map[string]interface{}, avoided []string) string {
	qry := "update " + table + " set "
	x := 0
	y := 0
	//* build our data field
	for key, val := range mymapdata {
		v := string(val.(string))
		if IsInArray(avoided, key) {
			if v != "" {
				if x == 0 {
					qry = qry + " " + key + "='" + v + "' "
					x++
				} else {
					qry = qry + ", " + key + "='" + v + "' "
					x++
				}
			}
		}
	}

	//* build our condition where clause

	for key, val := range mymapcond {
		v := string(val.(string))
		if v != "" {
			if y == 0 {
				qry = qry + " where " + key + "='" + v + "' "
				y++
			} else {
				qry = qry + " and " + key + "='" + v + "' "
				y++
			}
		}
	}

	return qry
}
func IsInArray(array []string, target string) bool {
	boo := true
	for _, b := range array {
		b = strings.ToLower(b)
		target = strings.ToLower(target)
		if b == target {
			boo = false;
			break
		}
	}
	return boo
}




