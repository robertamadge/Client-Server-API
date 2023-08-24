package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	gorm.Model
	ExchangeRateID int
	Name           string `json:"name"`
	Bid            string `json:"bid"`
}

type ExchangeRate struct {
	gorm.Model
	USDBRL USDBRL `json:"USDBRL" gorm:"foreignKey:ExchangeRateID"`
}

func main() {
	muxServer := http.NewServeMux()

	muxServer.HandleFunc("/cotacao", CotacaoHandler)

	err := http.ListenAndServe(":8080", muxServer)
	if err != nil {
		return
	}
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxServer := r.Context()

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Request processed with success.")
	case <-ctxServer.Done():
		log.Println("Request cancelled.")
	}

	req, err := http.NewRequestWithContext(ctxServer,
		"GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error in the response", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

	timeoutDB := 210 * time.Millisecond
	ctxDB, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	dsn := "database.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Close the database connection when done
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// AutoMigrate creates the tables if they don't exist
	err = db.AutoMigrate(&ExchangeRate{}, &USDBRL{})
	if err != nil {
		http.Error(w, "Error in migration", http.StatusInternalServerError)
		return
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	errDB := db.WithContext(ctxDB).Create(&exchangeRate).Error
	if errDB != nil {
		log.Println("Error inserting exchange rate:", errDB)
		log.Println("Request cancelled.")
	} else {
		log.Println("Exchange rate inserted successfully.")
		log.Println("Request processed with success.")
	}
}
