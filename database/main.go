package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bid  string `json:"bid"`
}

type ExchangeRate struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func main() {
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connString := "root:root@tcp(localhost:3306)/goexpert"
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
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
	newExchangeRate := NewExchangeRate(exchangeRate.USDBRL.Name, exchangeRate.USDBRL.Bid)

	select {
	case <-time.After(timeout):
		log.Println("Request processed with success.")
	case <-ctx.Done():
		log.Println("Request cancelled.")
	default:
		err := insertExchangeRate(db, *newExchangeRate)
		if err != nil {
			log.Println("Error inserting exchange rate:", err)
		} else {
			log.Println("Exchange rate inserted successfully.")
		}
	}
}

func NewExchangeRate(name, bid string) *ExchangeRate {
	return &ExchangeRate{
		USDBRL: USDBRL{
			ID:   uuid.New().String(),
			Name: name,
			Bid:  bid,
		},
	}
}

func insertExchangeRate(db *sql.DB, exchangeRate ExchangeRate) error {
	statment, err := db.Prepare("insert into exchange_rate(id, name, bid) values(?, ?, ?)")
	if err != nil {
		return err
	}

	defer statment.Close()

	_, err = statment.Exec(exchangeRate.USDBRL.ID, exchangeRate.USDBRL.Name, exchangeRate.USDBRL.Bid)
	if err != nil {
		return err
	}

	return nil
}
