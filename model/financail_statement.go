package model

type FinancialStatementInfo struct { // the cassandra corresponding table is [ finacialStatementCategory ]
	OrgCode string
	Id string
	Type string
	Category string
	SubCategory string
	StartCode int64
	EndCode int64
	SessionId string
	Date string
}
