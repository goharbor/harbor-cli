package api

type ListFlags struct {
	Name     string
	Page     int64
	PageSize int64
	Q        string
	Sort     string
	Public   bool
}
