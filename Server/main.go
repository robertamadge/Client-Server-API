package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	//db, err := sql.Open("sqlite3", "currency.db")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//
	//_, err = db.Exec(`
	//	CREATE TABLE IF NOT EXISTS currency (
	//		id INTEGER PRIMARY KEY,
	//		dollar FLOAT
	//	);
	//`)
	//if err != nil {
	//	log.Fatal(err)
	//}

	muxServer := http.NewServeMux()

	muxServer.HandleFunc("/cotacao", CotacaoHandler)

	http.ListenAndServe(":8080", muxServer)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Request processed with success:server.")
	case <-ctx.Done():
		log.Println("Request cancelled:server.")
	}

	req, err := http.NewRequestWithContext(ctx,
		"GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error in the response server", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}
	//var usdBrl USDBRL
	//err = json.Unmarshal(body, &usdBrl)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(body)

	w.Header().Set("Content-Type", "application/json")

	w.Write(body)
}
