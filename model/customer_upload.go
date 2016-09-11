package model

type CustomerUpload struct {
	UpRef string
	OrgCode string
	Id string
	Refrence string
	Date string
	AccountingCode string
	Year int
	Month int
	Day int
	AccounttingSystem string
	TxnType string
	DebitValue float64
	CreditValue float64
	EntryDescription string
	EntryCategory string
	EntrySubCategory string
	CsvStringInput string
	MappingCode string
}

