package response

type LineResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type PieResponse struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}
