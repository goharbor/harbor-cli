package api

type ListFlags struct {
	Name     string
	Page     int64
	PageSize int64
	Q        string
	Sort     string
	Public   bool
}

type ListQuotaFlags struct {
	PageSize    int64
	Page        int64
	Sort        string
	Reference   string
	ReferenceID string
}
