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
	"strings"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/offline"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/spf13/cobra"
)

var code string
var historyDays int

func init() {
	plotStockChartCmd.Flags().StringVarP(&code, "code", "c", "AAPL", "the stock's symbol")
	plotStockChartCmd.Flags().IntVarP(&historyDays, "days", "d", 30, "history days")
	rootCmd.AddCommand(plotStockChartCmd)
}

var plotStockChartCmd = &cobra.Command{
	Use:   "plot_stock_chart",
	Short: "draw stock chart",
	Long:  `draw stock chart`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./stock_history")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		querySQL := fmt.Sprintf("SELECT * FROM stock_daily_his WHERE symbol = '%s' ORDER BY date DESC LIMIT %d", code, historyDays)
		rows, err := db.Query(querySQL)
		if err != nil {
			log.Fatalf("Failed to query data: %v", err)
		}
		defer rows.Close()

		dates := []string{}
		closePrices := []float64{}

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
			dates = append(dates, strings.Split(date, "T")[0])
			closePrices = append(closePrices, close)
		}

		fig := &grob.Fig{
			Data: []types.Trace{
				&grob.Scatter{
					X: types.DataArray(dates),
					Y: types.DataArray(closePrices),
				},
			},
			Layout: &grob.Layout{
				Title: &grob.LayoutTitle{
					Text: types.StringType(fmt.Sprintf("A Figure of %s stock history", code)),
				},
			},
		}

		offline.Show(fig)
	},
}
