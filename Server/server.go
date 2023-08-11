package Server

import (
	"fmt"
	"io"
	"net/http"
)

// Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida
// , sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo
// para conseguir persistir os dados no banco deverá ser de 10ms.
type Cambio struct {
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

func RunServer() {
	muxServer := http.NewServeMux()

	cambio := Cambio{}
	muxServer.HandleFunc("/cotacao", cambio.CotacaoHandler)

	http.ListenAndServe(":8081", muxServer)
}

func (c Cambio) CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		panic(err)
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")

	fmt.Println(string(body))
	w.Write(body)
}
