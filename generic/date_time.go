package generic

import (
	"time"
	"fmt"
	"log"
	"strings"
)

func StringToTime(date string,timein string)time.Time{
	p := fmt.Println
	//t := time.Now()
	//p(t.Format(time.RFC3339))

	myformat :=date+"T"+timein+"+00:00"
	t1, e := time.Parse(
		time.RFC3339,
		myformat)

	if e !=nil{
		log.Println("ERROR DATE PARSE > ",e)

	}
	p("XXXXX>> ",t1)
t :=t1
	fmt.Printf(":>> %d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return t1
}

func GetDurationUpToNow(dateFrom,timeFrom string)time.Duration{

	loc, _ := time.LoadLocation("UTC")

	startDt :=StringToTime(dateFrom,timeFrom).UTC()

	createdAt := startDt.In(loc).Add(-2 * time.Hour)

	currDate :=time.Now().UTC()

	startDt.Sub(createdAt)


	duration :=startDt.Sub(currDate)// time.Since(startDt)
	//fmt.Println(duration.Hours())
	return duration
}

func GetDateAndTimeString_AddDays(numDay int) (string, string) {
	mydate := time.Now()
	threeDays := time.Hour * 24 * time.Duration(numDay)
	now := time.Now()
	mydate = now.Add(threeDays)
	arr := strings.Split(fmt.Sprintln(mydate.Format("2006-01-02 15:04:05")), " ")
	date := arr[0]
	time := arr[1]
	return strings.TrimSpace(date), strings.TrimSpace(time)
}
func GetDateAndTimeString() (string, string) {
	mydate := time.Now()
	arr := strings.Split(fmt.Sprintln(mydate.Format("2006-01-02 15:04:05")), " ")
	date := arr[0]
	time := arr[1]
	return strings.TrimSpace(date), strings.TrimSpace(time)
}

func GetDurationValueIntoString(duration time.Duration)(string,string,string,string){
	hour :=""
	minute :=""
	second :=""
	millisecond :=""
	hour = fmt.Sprintf("%.0f",duration.Hours())
	minute = fmt.Sprintf("%.0f",duration.Minutes())
	second = fmt.Sprintf("%.0f",duration.Seconds())
	millisecond =fmt.Sprintf("%d",duration.Nanoseconds())
	return hour,minute,second,millisecond
}


func GetDurationDatetimeUpToNow(date1 string) (string, time.Duration) {
	dt1 := ConvertStringTimeGolangFormateToTime(date1)
	dt2 := ConvertStringTimeGolangFormateToTime(time.Now().String())
	myduration := dt2.Sub(dt1)
	return myduration.String(), myduration
}
func GetDurationDateToDate(date1 string, date2 string) (string, time.Duration) {
	dt1 := ConvertStringTimeGolangFormateToTime(date1)
	dt2 := ConvertStringTimeGolangFormateToTime(date2)
	myduration := dt2.Sub(dt1)
	return myduration.String(), myduration
}

func ConvertStringTimeGolangFormateToTime(str string) time.Time {
	//2016-04-23T14:19:43.064+0200
	split1 := strings.Split(str, ".")
	dtstr := strings.Replace(split1[0], "T", " ", 1)
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dtstr)
	if err != nil {
		fmt.Println(err)
	}
	return t
}