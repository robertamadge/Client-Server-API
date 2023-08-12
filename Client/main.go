package Client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Request processed with success.")
	case <-ctx.Done():
		log.Println("Request cancelled.")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error in the request client:", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error in the response client:", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body client:", err)
		return
	}

	//var usdBrl USDBRL
	//err = json.Unmarshal(body, &usdBrl)
	//if err != nil {
	//	panic(err)
	//}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bid := string(body)
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		panic(err)
	}
}
