package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	cotacao, err := BuscaCotacao()

	// Se der algum erro, escreve no log do servidor e retorna um erro 500 para quem chamou a API
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Inicia o BD e cria a tabela de cotação
	db, err := initBD()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer db.Close()

	// Insere a cotação na tabela
	err = insertCotacao(db, cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Escrever no header do respose e retorna apenas o valor da cotação cotacao.USDBRL.Bid
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao.USDBRL.Bid)
}

func BuscaCotacao() (*Cotacao, error) {
	// Cria o contexto com o time limit de 200ms
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	// Cria a requisição
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func initBD() (*sql.DB, error) {
	// Valida se o arquivo database.db está criado e se não tiver, cria ele
	_, err := os.OpenFile("./database.db", os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Carrega o BD
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			codein TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid TEXT,
			ask TEXT,
			timestamp TEXT,
			create_date TEXT
		);
	`)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertCotacao(db *sql.DB, cotacao *Cotacao) error {
	// Cria o contexto com o time limit de 10ms
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(
		ctx,
		`INSERT INTO cotacoes (
			code, codein, name,
			high, low, var_bid,
			pct_change, bid, ask,
			timestamp, create_date
		) VALUES (
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?
		)`,
		cotacao.USDBRL.Code,
		cotacao.USDBRL.Codein,
		cotacao.USDBRL.Name,
		cotacao.USDBRL.High,
		cotacao.USDBRL.Low,
		cotacao.USDBRL.VarBid,
		cotacao.USDBRL.PctChange,
		cotacao.USDBRL.Bid,
		cotacao.USDBRL.Ask,
		cotacao.USDBRL.Timestamp,
		cotacao.USDBRL.CreateDate,
	)

	if err != nil {
		return err
	}

	return nil
}
