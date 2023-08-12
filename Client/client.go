package Client

import (
	"context"
	"encoding/json"
	"fmt"
	server "github.com/robertamadge/Client-Server-API/Server"
	"io"
	"net/http"
	"os"
	"time"
)

//Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var cambio server.Cotacao
	err = json.Unmarshal(body, &cambio)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bid := cambio.USDBRL.Bid
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		panic(err)
	}
}
