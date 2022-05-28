package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Message struct {
	Title   string
	Content string
	Author  string
}

func main() {
	client := setClient()

	ctx := context.Background() // 執行ES請求需要提供一個上下文對象

	// 定義資料轉[]byte
	data, _ := json.Marshal(
		Message{
			Title:   "Test Title",
			Content: "Test Test",
			Author:  "Pearce",
		},
	)

	// 使用client創建一個新的資料
	req := esapi.IndexRequest{
		Index:      "test",
		DocumentID: strconv.Itoa(1),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// 執行指令
	res, err := req.Do(ctx, client)
	if err != nil {
		log.Fatalf("Error getting response: %s\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d\n", res.Status(), 1)
	} else {
		// 將資料反序列化
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s\n", err)
		} else {
			// 印出狀態和版本
			log.Printf("[%s] %s; version=%d\n", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

}

// Client 設定 詳細資料請看elasticsearch.go
func setClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
			"https://localhost:9201",
		}, // Service ip
		Username: "user",     // Authentication Username
		Password: "password", // Authentication password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s\n", err)
		fmt.Printf("連接失敗: %s\n", err)
	} else {
		fmt.Printf("連接成功\n")
	}

	return client
}
