package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "/tmp/stock_history.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS stock_daily_his (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol VARCHAR NOT NULL,
		date DATE,
		open NUMERIC,
		high NUMERIC,
		low  NUMERIC,
		close NUMERIC,
		volume NUMERIC
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	createTableSQL = `CREATE TABLE IF NOT EXISTS console_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATE,
		context VARCHAR NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	createTableSQL = `CREATE TABLE IF NOT EXISTS fund_daily_his (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fund_code VARCHAR NOT NULL,
		fund_name VARCHAR NOT NULL,
		date DATE,
		avg_net_worth NUMERIC,
		accumulated_net_worth NUMERIC,
		daily_growth_rate NUMERIC
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	return db
}

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		db: InitDB(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetConsoleLogs() string {
	console_log := ""
	querySQL := "select * from console_log order by id desc limit 10;"
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	for rows.Next() {
		var id int
		var time string
		var context string

		err = rows.Scan(&id, &time, &context)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		console_log += fmt.Sprintf("id: %d time: %s context: %s \n", id, time, context)
	}
	return console_log
}

func (a *App) SyncFundHisData(fund_code string) {
	fundHisData := FetchFundData(fund_code)
	for _, fundData := range fundHisData {
		deleteSQL := `DELETE FROM fund_daily_his WHERE fund_code = ? and date = ?`
		_, err := a.db.Exec(deleteSQL, fundData.FundCode, fundData.DateTime)
		if err != nil {
			log.Fatalf("Failed to delete data: %v", err)
		}
		insertSQL := `INSERT INTO fund_daily_his (fund_code, fund_name, date, avg_net_worth, accumulated_net_worth, daily_growth_rate) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = a.db.Exec(insertSQL, fundData.FundCode, fundData.FundName, fundData.DateTime, fundData.UnitNetValue, fundData.CumulativeNetValue, fundData.DailyGrowthRate)
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}
		insertSQL = `INSERT INTO console_log (time, context) VALUES (?, ?)`
		_, err = a.db.Exec(insertSQL, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("sync fund data success, code %s, date %s", fundData.FundCode, fundData.DateTime))
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}
	}
}

func (a *App) SyncStockHisData(stock_code string, token string) {
	stockHisData := FetchStockDaliyHistoryFromRemote(stock_code, token)
	for date, data := range stockHisData.TimeSeries {
		fmt.Printf("  Date: %s\n", date)
		fmt.Printf("    Open: %s\n", data.Open)
		fmt.Printf("    High: %s\n", data.High)
		fmt.Printf("    Low: %s\n", data.Low)
		fmt.Printf("    Close: %s\n", data.Close)
		fmt.Printf("    Volume: %s\n", data.Volume)

		deleteSQL := `DELETE FROM stock_daily_his WHERE symbol = ? and date = ?`
		_, err := a.db.Exec(deleteSQL, stock_code, date)
		if err != nil {
			log.Fatalf("Failed to delete data: %v", err)
		}

		insertSQL := `INSERT INTO stock_daily_his (symbol, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = a.db.Exec(insertSQL, stock_code, date, data.Open, data.High, data.Low, data.Close, data.Volume)
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}

		insertSQL = `INSERT INTO console_log (time, context) VALUES (?, ?)`
		_, err = a.db.Exec(insertSQL, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("sync stock data success, code %s, date %s", stock_code, date))
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}
	}
}

func (a *App) GetChartDataLabels(stock_code string, start_date string, end_date string) []string {
	var labels []string
	querySQL := fmt.Sprintf("SELECT * FROM stock_daily_his WHERE symbol = '%s' AND date BETWEEN '%s' AND '%s' ORDER BY date ASC;", stock_code, start_date, end_date)
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	for rows.Next() {
		var id int
		var fund_code string
		var fund_name string
		var date string
		var avg_net_worth float64
		var accumulated_net_worth float64
		var daily_growth_rate float64

		err = rows.Scan(&id, &fund_code, &fund_name, &date, &avg_net_worth, &accumulated_net_worth, &daily_growth_rate)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("id %d, fund_code %s, date %s, accumulated net worth %f \n", id, fund_code, strings.Split(date, "T")[0], accumulated_net_worth)
		labels = append(labels, strings.Split(date, "T")[0])
	}
	defer rows.Close()
	fmt.Println(labels)
	return labels
}

func (a *App) GetFundName(fund_code string) string {
	querySQL := fmt.Sprintf("SELECT * FROM fund_daily_his WHERE fund_code = '%s' limit 1;", fund_code)
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	fundName := ""
	for rows.Next() {
		var id int
		var fund_code string
		var fund_name string
		var date string
		var avg_net_worth float64
		var accumulated_net_worth float64
		var daily_growth_rate float64

		err = rows.Scan(&id, &fund_code, &fund_name, &date, &avg_net_worth, &accumulated_net_worth, &daily_growth_rate)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fundName = fund_name
	}
	defer rows.Close()
	return fundName
}

func (a *App) GetFundChartDataLabels(fund_code string, start_date string, end_date string) []string {
	var labels []string
	querySQL := fmt.Sprintf("SELECT * FROM fund_daily_his WHERE fund_code = '%s' AND date BETWEEN '%s' AND '%s' ORDER BY date ASC;", fund_code, start_date, end_date)
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	for rows.Next() {
		var id int
		var fund_code string
		var fund_name string
		var date string
		var avg_net_worth float64
		var accumulated_net_worth float64
		var daily_growth_rate float64

		err = rows.Scan(&id, &fund_code, &fund_name, &date, &avg_net_worth, &accumulated_net_worth, &daily_growth_rate)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("id %d, fund_code %s, date %s, accumulated net worth %f \n", id, fund_code, strings.Split(date, "T")[0], accumulated_net_worth)
		labels = append(labels, strings.Split(date, "T")[0])
	}
	defer rows.Close()
	fmt.Println(labels)
	return labels
}

func (a *App) GetFundChartData(fund_code string, start_date string, end_date string) [][]float64 {
	const COLS = 3

	var matrix [][]float64

	querySQL := fmt.Sprintf("SELECT * FROM fund_daily_his WHERE fund_code = '%s' AND date BETWEEN '%s' AND '%s' ORDER BY date ASC;", fund_code, start_date, end_date)
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	for rows.Next() {
		var id int
		var fund_code string
		var fund_name string
		var date string
		var avg_net_worth float64
		var accumulated_net_worth float64
		var daily_growth_rate float64

		err = rows.Scan(&id, &fund_code, &fund_name, &date, &avg_net_worth, &accumulated_net_worth, &daily_growth_rate)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("id %d, fund_code %s, date %s, accumulated net worth %f \n", id, fund_code, strings.Split(date, "T")[0], accumulated_net_worth)
		rowData := make([]float64, COLS)
		rowData[0] = avg_net_worth
		rowData[1] = accumulated_net_worth
		rowData[2] = daily_growth_rate
		matrix = append(matrix, rowData)
	}
	defer rows.Close()
	return matrix
}

func (a *App) GetChartData(stock_code string, start_date string, end_date string) [][]float64 {

	const COLS = 4

	var matrix [][]float64

	querySQL := fmt.Sprintf("SELECT * FROM stock_daily_his WHERE symbol = '%s' AND date BETWEEN '%s' AND '%s' ORDER BY date ASC;", stock_code, start_date, end_date)
	rows, err := a.db.Query(querySQL)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}

	for rows.Next() {
		var id int
		var symbol string
		var date string
		var open float64
		var high float64
		var low float64
		var close float64
		var volume float64
		err = rows.Scan(&id, &symbol, &date, &open, &high, &low, &close, &volume)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("id %d, symbol %s, date %s, close %f \n", id, symbol, strings.Split(date, "T")[0], close)

		rowData := make([]float64, COLS)
		rowData[0] = open
		rowData[1] = close
		rowData[2] = low
		rowData[3] = high
		matrix = append(matrix, rowData)
	}

	defer rows.Close()
	return matrix
}
