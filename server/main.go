package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyQuotation struct {
	Coin Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
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
	http.HandleFunc("/all", handlerGetAllQuotation)
	http.HandleFunc("/", handlerGetQuotation)
	http.ListenAndServe(":8080", nil)
}

func handlerGetAllQuotation(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "databsase.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer database.Close()

	quotations, err := selectAllQuotation(database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var coins []*Usdbrl
	for _, quotation := range quotations {
		coins = append(coins, &quotation.Coin)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coins)
}

func selectAllQuotation(db *sql.DB) ([]CurrencyQuotation, error) {

	rows, err := db.Query("select code, codein, name from quotation")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotations []CurrencyQuotation
	for rows.Next() {
		var coin Usdbrl
		err = rows.Scan(&coin.Code, &coin.Codein, &coin.Name)
		if err != nil {
			return nil, err
		}

		quotation := CurrencyQuotation{Coin: coin}
		quotations = append(quotations, quotation)
	}

	return quotations, nil
}

func handlerGetQuotation(w http.ResponseWriter, r *http.Request) {
	quotation, err := getQuotation(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(quotation)
	err = createDatase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	database, err := sql.Open("sqlite3", "databsase.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer database.Close()

	err = createTable(database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = insertQuotation(database, quotation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getQuotation(w http.ResponseWriter, r *http.Request) (*CurrencyQuotation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, err
	}

	var quotation CurrencyQuotation
	error = json.Unmarshal(body, &quotation)
	if error != nil {
		return nil, err
	}

	return &quotation, nil
}

func createDatase() error {
	f, err := os.Create("cotacao.db")
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func createTable(db *sql.DB) error {
	quotation_table := `CREATE TABLE IF NOT EXISTS quotation (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "code" TEXT,
        "codein" TEXT,
        "name" TEXT,
		"high" TEXT,
		"low" TEXT,
		"varBid" TEXT,
		"pctChange" TEXT,
		"bid" TEXT,
		"ask" TEXT,
		"timestamp" TEXT,
		"create_date" TEXT);`
	query, err := db.Prepare(quotation_table)
	if err != nil {
		return err
	}
	query.Exec()
	fmt.Println("Table created successfully!")

	return nil
}

func insertQuotation(db *sql.DB, quotation *CurrencyQuotation) error {

	stmt, err := db.Prepare("insert into quotation(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(quotation.Coin.Code, quotation.Coin.Codein, quotation.Coin.Name, quotation.Coin.High, quotation.Coin.Low, quotation.Coin.VarBid, quotation.Coin.PctChange, quotation.Coin.Bid, quotation.Coin.Ask, quotation.Coin.Timestamp, quotation.Coin.CreateDate)
	if err != nil {
		return err
	}

	return nil
}
