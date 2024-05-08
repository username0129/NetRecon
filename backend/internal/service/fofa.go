package service

import (
	"backend/internal/config"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type FofaService struct{}

var (
	FofaServiceApp = new(FofaService)
)

// FofaResponse 用于解析 Fofa API 返回的 JSON 数据
type FofaResponse struct {
	Error   bool       `json:"error"`
	ErrMsg  string     `json:"errmsg"`
	Results [][]string `json:"results"`
	Mode    string     `json:"mode"`
	Page    int64      `json:"page"`
	Query   string     `json:"query"`
	Size    int64      `json:"size"`
}

// BaseEncode Base64 加密接口
func BaseEncode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (fs *FofaService) FofaSearch(fofa config.Fofa, req request.FofaSearchRequest) (fofalist response.FofaSearchResponse, err error) {
	baseURL := fofa.Url + "api/v1/search/all"
	data := url.Values{}
	data.Set("email", fofa.Mail)
	data.Set("key", fofa.Key)
	data.Set("qbase64", url.QueryEscape(BaseEncode(req.Query)))
	data.Set("page", strconv.Itoa(req.Page))
	data.Set("size", strconv.Itoa(req.PageSize))
	data.Set("fields", "host,title,ip,domain,port,protocol,country_name,region,city,icp")

	// 发起 HTTP GET 请求
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, data.Encode()))
	if err != nil {
		return response.FofaSearchResponse{}, errors.New("发送 FOFA 查询请求失败")
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.FofaSearchResponse{}, errors.New("读取 FOFA 响应失败")
	}

	// 解析 JSON 数据
	var fofaResp FofaResponse
	if err := json.Unmarshal(body, &fofaResp); err != nil {
		return response.FofaSearchResponse{}, errors.New("JSON 反序列化失败")
	}
	fofalist.Total = fofaResp.Size
	if fofaResp.Error {
		fofalist.Status = false
		fofalist.Message = fofaResp.ErrMsg
	} else {
		if len(fofaResp.Results) == 0 {
			fofalist.Message = "未查询到有效数据"
		} else {
			var temp string
			fofalist.Status = true
			for _, result := range fofaResp.Results {
				// 拼接 HTTP
				if !strings.Contains(result[0], "://") && (result[5] == "http" || result[5] == "https") {
					temp = result[5] + "://" + result[0]
				} else {
					temp = result[0]
				}
				fofalist.Results = append(fofalist.Results, response.FofaResponse{
					URL:      temp,
					Title:    result[1],
					IP:       result[2],
					Domain:   result[3],
					Port:     result[4],
					Protocol: result[5],
					Country:  result[6],
					Region:   result[7],
					City:     result[8],
					ICP:      result[9],
				})
			}
		}
	}
	return fofalist, nil
}
