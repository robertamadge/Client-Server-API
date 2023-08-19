package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type USDBRL struct {
	gorm.Model
	ExchangeRateID int
	Name           string `json:"name"`
	Bid            string `json:"bid"`
}

type ExchangeRate struct {
	gorm.Model
	USDBRL USDBRL `gorm:"foreignKey:ExchangeRateID"`
}

func main() {
	timeout := 10 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dsn := "db/database.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	db.AutoMigrate(&ExchangeRate{}, &USDBRL{})

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

	select {
	case <-time.After(timeout):
		log.Println("Request processed with success.")
	case <-ctx.Done():
		log.Println("Request cancelled.")
	default:
		errDB := db.Create(&exchangeRate).Error
		if errDB != nil {
			log.Println("Error inserting exchange rate:", errDB)
		} else {
			log.Println("Exchange rate inserted successfully.")
		}
	}
}
