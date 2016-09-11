package model

type UploadSetting struct {
	OrgCode string
	SeesionId string
	Id string
	MapingType string
	DateFormat string
	AccountSystem string
	CodeColom int
	DescColom int
	DebitColom int
	CreditColom int
	DateColom int
	RowStart int
	Status string
}

func(upSetting *UploadSetting)NewSetting(){

}
func(upSetting *UploadSetting)GetDefault()UploadSetting{
	return  DefautSetting()
}

func DefautSetting() UploadSetting{
	upset := UploadSetting{}
	upset.OrgCode = "MM01"
	upset.Id = "MM2016-08-19x16-04"
	upset.MapingType = "MT001" // MT001 -> mapping range
	upset.DateFormat = "dd/mm/yyyy"
	upset.AccountSystem = "pastel"
	upset.CodeColom = 2
	upset.DescColom = 3
	upset.DebitColom = 4
	upset.CreditColom = 5
	upset.DateColom = 2
	upset.RowStart = 4
	upset.Status = "active"
	return upset
}


