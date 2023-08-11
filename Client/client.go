package Client

import (
	"encoding/json"
	"fmt"
	server "github.com/robertamadge/Client-Server-API/Server"
	"io"
	"net/http"
	"os"
)

//Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

func RunClient() {
	client := http.Client{}

	resp, err := client.Get("http://localhost:8081/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var cambio server.Cambio
	err = json.Unmarshal(body, &cambio)
	if err != nil {
		panic(err)
	}

	bid := cambio.USDBRL.Bid
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		panic(err)
	}
}
