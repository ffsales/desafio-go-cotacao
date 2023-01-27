package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	http.HandleFunc("/", handlerGetQuotation)
	http.ListenAndServe(":8080", nil)
}

func handlerGetQuotation(w http.ResponseWriter, r *http.Request) {
	// quotation, err := getQuotation(w, r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Println(quotation)
	createDatase()
	database, err := sql.Open("sqlite3", "databsase.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer database.Close()
	createTable(database)
	// connectDatabase()
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

func createDatase() {
	f, err := os.Create("cotacao.db")
	if err != nil {
		panic(err)
	}
	f.Close()
}

func createTable(db *sql.DB) {
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
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table created successfully!")
}

func connectDatabase() {

	//create table if not exist

	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM primeiro_teste;")
	if err != nil {
		panic(err)
	}

	var teste_id int
	var teste_name string

	for rows.Next() {
		err = rows.Scan(&teste_id, &teste_name)
		if err != nil {
			panic(err)
		}

		fmt.Println(teste_id)
		fmt.Println(teste_name)
	}
	rows.Close()

}
