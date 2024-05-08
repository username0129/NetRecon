package request

type FofaSearchRequest struct {
	Query    string `json:"query"`
	PageSize int    `json:"pageSize"`
	Page     int    `json:"page"`
}
