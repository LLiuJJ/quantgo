// Copyright 2025 The quantgo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"database/sql"
	"fmt"
	"log"

	quantgo_crawler "github.com/LLiuJJ/quantgo/crawler"
	"github.com/spf13/cobra"
)

var stockCode string
var apiToken string

func init() {
	syncStockDataCmd.Flags().StringVarP(&stockCode, "code", "c", "AAPL", "the stock's symbol")
	syncStockDataCmd.Flags().StringVarP(&apiToken, "token", "t", "xxxx", "the api token, you can get free api token from https://www.alphavantage.co/support/#api-key")
	rootCmd.AddCommand(syncStockDataCmd)
}

var syncStockDataCmd = &cobra.Command{
	Use:   "sync_stock_data",
	Short: "Sync stock history data from remote",
	Long:  `Sync stock history data from remote`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stock code", stockCode)
		fmt.Println("api token", apiToken)

		db, err := sql.Open("sqlite3", "./stock_history")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()
		stockHisData := quantgo_crawler.FetchStockDaliyHistoryFromRemote(stockCode, apiToken)
		for date, data := range stockHisData.TimeSeries {
			fmt.Printf("  Date: %s\n", date)
			fmt.Printf("    Open: %s\n", data.Open)
			fmt.Printf("    High: %s\n", data.High)
			fmt.Printf("    Low: %s\n", data.Low)
			fmt.Printf("    Close: %s\n", data.Close)
			fmt.Printf("    Volume: %s\n", data.Volume)

			deleteSQL := `DELETE FROM stock_daily_his WHERE symbol = ? and date = ?`
			_, err = db.Exec(deleteSQL, stockCode, date)
			if err != nil {
				log.Fatalf("Failed to delete data: %v", err)
			}

			insertSQL := `INSERT INTO stock_daily_his (symbol, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)`
			_, err = db.Exec(insertSQL, stockCode, date, data.Open, data.High, data.Low, data.Close, data.Volume)
			if err != nil {
				log.Fatalf("Failed to insert data: %v", err)
			}
		}
	},
}
