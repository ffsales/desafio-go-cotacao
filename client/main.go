package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type QuotationResponse struct {
	Bid string `json:"bid"`
}

func main() {

	quotation, err := getQuotationServer()

	if err != nil {
		panic(err)
	}

	fmt.Println(quotation)
}

func getQuotationServer() (*QuotationResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var quotation QuotationResponse

	error = json.Unmarshal(body, &quotation)
	if error != nil {
		return nil, error
	}

	return &quotation, nil
}
