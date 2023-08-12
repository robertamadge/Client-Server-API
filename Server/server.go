package Server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Cotacao struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	muxServer := http.NewServeMux()

	cotacao := Cotacao{}
	muxServer.HandleFunc("/cotacao", cotacao.CotacaoHandler)

	http.ListenAndServe(":8080", muxServer)
}

func (c Cotacao) CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	select {
	case <-time.After(5 * time.Second):
		log.Println("Request processada com sucesso.")
		w.Write([]byte("Request processada com sucesso"))
	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
		http.Error(w, "Request cancelada pelo cliente", http.StatusRequestTimeout)
	}

	req, err := http.NewRequestWithContext(ctx,
		"GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)
	if err != nil {
		fmt.Println("Error in the request server:", err)
		return
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading body server:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println(string(body))
	w.Write(body)
}
