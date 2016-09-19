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

func (cup *CustomerUpload)DetermineTransctionType( )string{

	if cup.CreditValue !=0 && cup.DebitValue !=0{
		cup.TxnType = "BOTH"
		return "BOTH"
	}

	if cup.CreditValue ==0 && cup.DebitValue ==0{
		cup.TxnType = "FUTURE"
		return "FUTURE"
	}

	if cup.CreditValue ==0 && cup.DebitValue !=0{
		cup.TxnType = "DEBIT"
		return "DEBIT"
	}

	if cup.CreditValue !=0 && cup.DebitValue ==0{
		cup.TxnType = "CREDIT"
		return "CREDIT"
	}

	return "UNKNOW"
}

