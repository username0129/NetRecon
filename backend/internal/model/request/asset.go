package request

type AddAssetRequest struct {
	Domains string `json:"domains"`
	IPs     string `json:"ips"`
}
