package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Struct que mapeia o retorno da API de cotação
type Cotacao struct {
	USDBRL struct {
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
	} `json:"USDBRL"`
}

func main() {
	// Cria o servidor na porta 8080 e cria a rota /cotacao que chamará o hander BuscaCotacao
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	// Valida que a rota chamada foi a /cotacao
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Busca a cotação
	cotacao, error := BuscaCotacao()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Escrever no header do respose e retorna
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func BuscaCotacao() (*Cotacao, error) {
	// Cria o contexto com o time limit de 200ms
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	// Cria a requisição
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	// Faz a requisição para recuperar a cotação
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Recupera o body retornado
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	// Faz o json encode e mapeia na struct Cotacao
	var cotacao Cotacao
	error = json.Unmarshal(body, &cotacao)
	if error != nil {
		return nil, error
	}

	return &cotacao, nil
}
