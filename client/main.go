package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	cotacao, err := buscaCotacao()
	if err != nil {
		log.Println(err)
		return
	}

	err = escreveNoArquivo(cotacao)
	if err != nil {
		log.Println(err)
		return
	}
}

func buscaCotacao() (*CotacaoResponse, error) {
	// Cria o contexto com o time limit de 300ms
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	// Cria a requisição
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	// Faz a requisição para recuperar a cotação
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Recupera o body retornado
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Faz o json encode e mapeia na struct Cotacao
	var cotacao CotacaoResponse
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func escreveNoArquivo(cotacao *CotacaoResponse) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.Bid))

	if err != nil {
		return err
	}

	return nil
}
