package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

type ExchangeRate struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func main() {
	ctxClient, cancel := context.WithTimeout(context.Background(), 510*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxClient, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error in the request:", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error in the response:", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bid := exchangeRate.USDBRL.Bid
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		panic(err)
	}
}
