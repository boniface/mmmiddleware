package dbconn

import (
	"testing"
	"fmt"
	"github.com/gocql/gocql"
)


func TestCreateCassandraSession2(t *testing.T){
	CreateCassandraSession2()

}
func TestGetServerCassandra(t *testing.T) {
	p :=fmt.Println
	SetDbServerFilename()
	list :=GetServerCassandra("test")
	p("==>List server ",len(list))
	for x,ser:=range list{
		p(x," --> ",ser)
	}
}

func TestInitCreateDatabaseTablesSCHEMA(t *testing.T) {

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 0
	cluster.Keyspace  ="ckiosk"
	inner_session, errCass := cluster.CreateSession()
	//defer inner_session.Close()
	if errCass != nil {
		fmt.Println("errCass > ", errCass)
		//return errCass
	}
	//fmt.Println("CONNECTION TO CASSANDRA 100%")
	CassSession = inner_session


	session := CassSession


	qry :=`
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

