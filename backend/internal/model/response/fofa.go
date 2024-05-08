package response

type FofaSearchResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Total   int64          `json:"total"`
	Results []FofaResponse `json:"results"`
}

type FofaResponse struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Domain   string `json:"domain"`
	Protocol string `json:"protocol"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	City     string `json:"city"`
	ICP      string `json:"icp"`
}
