package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 参数
// q={query} 搜索关键词
// page={page} 页
// pagesize={pagesize} 每页数量

type BaiduResponse struct {
	Data struct {
		Documents struct {
			Data []struct {
				TechDocDigest struct {
					URL     string `json:"url"`
					Title   string `json:"title"`
					Summary string `json:"summary"`
				} `json:"techDocDigest"`
			} `json:"data"`
		} `json:"documents"`
	} `json:"data"`
}

func baiduKaifaSearch(query string, page int, pageSize int) ([]map[string]interface{}, error) {
	if query == "" || page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	url := fmt.Sprintf("https://kaifa.baidu.com/rest/v1/search?wd=%s&paramList=%d&page_size=%d", query, page, pageSize)

	headers := map[string]string{
		"Host":            "kaifa.baidu.com",
		"Connection":      "keep-alive",
		"Accept":          "application/json, text/plain, */*",
		"User-Agent":      "Mozilla/4.5 (compatible; HTTrack 3.0x; Windows 98)",
		"Referer":         "https://kaifa.baidu.com/",
		"Accept-Encoding": "", // 重要; 禁用 gzip 压缩
		"Accept-Language": "zh-CN,zh;q=0.9",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response BaiduResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range response.Data.Documents.Data {
		results = append(results, map[string]interface{}{
			"url":     doc.TechDocDigest.URL,
			"title":   doc.TechDocDigest.Title,
			"summary": doc.TechDocDigest.Summary,
		})
	}

	return results, nil
}

func Search(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "10"))

	results, err := baiduKaifaSearch(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}
