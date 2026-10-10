package main

import (
	"encoding/json"
	"io"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func BuscaCotacao() (*Cotacao, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c Cotacao
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
