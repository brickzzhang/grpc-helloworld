package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func BenchmarkSayHello(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		data := make(map[string]interface{})
		data["message"] = "test"
		bytesData, _ := json.Marshal(data)
		req, _ := http.NewRequest(
			"POST", "http://localhost:8081/v1/helloworld", bytes.NewReader(bytesData),
		)
		resp, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
