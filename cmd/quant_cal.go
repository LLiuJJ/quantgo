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
	"strconv"
	"strings"

	quantgo_crawler "github.com/LLiuJJ/quantgo/crawler"
	quantgo_algo "github.com/LLiuJJ/quantgo/tradingalgo"
	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/offline"
	"github.com/MetalBlueberry/go-plotly/pkg/types"

	"github.com/spf13/cobra"
)

var quantStockcode string
var quantStockHistoryDays int
var quantAlgo string

func init() {
	quantCalCmd.Flags().StringVarP(&quantAlgo, "algo", "a", "turtle", "the Trading Algorithm")
	quantCalCmd.Flags().StringVarP(&quantStockcode, "code", "c", "AAPL", "the stock's symbol")
	quantCalCmd.Flags().IntVarP(&quantStockHistoryDays, "days", "d", 30, "history days")
	rootCmd.AddCommand(quantCalCmd)
}

var quantCalCmd = &cobra.Command{
	Use:   "quant_cal",
	Short: "draw stock chart",
	Long:  `draw stock chart`,
	Run: func(cmd *cobra.Command, args []string) {
		switch quantAlgo {
		case "turtle":
			{
				db, err := sql.Open("sqlite3", "./stock_history")
				if err != nil {
					log.Fatalf("Failed to open database: %v", err)
				}
				defer db.Close()

				querySQL := fmt.Sprintf("SELECT * FROM stock_daily_his WHERE symbol = '%s' ORDER BY date DESC LIMIT %d", quantStockcode, quantStockHistoryDays)
				rows, err := db.Query(querySQL)
				if err != nil {
					log.Fatalf("Failed to query data: %v", err)
				}
				defer rows.Close()

				priceDatas := []quantgo_crawler.DailyData{}
				dates := []string{}
				closePrices := []float64{}

				for rows.Next() {
					dailyData := quantgo_crawler.DailyData{}
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
					dailyData.Close = strconv.FormatFloat(close, 'f', 2, 64)
					dailyData.High = strconv.FormatFloat(high, 'f', 2, 64)
					dailyData.Low = strconv.FormatFloat(low, 'f', 2, 64)
					dailyData.Open = strconv.FormatFloat(open, 'f', 2, 64)
					dailyData.Volume = strconv.FormatFloat(volume, 'f', 2, 64)
					priceDatas = append([]quantgo_crawler.DailyData{dailyData}, priceDatas...)
					closePrices = append([]float64{close}, closePrices...)
					dates = append([]string{strings.Split(date, "T")[0]}, dates...)
					fmt.Printf("id %d, symbol %s, date %s, close %f , high %f, low %f \n", id, symbol, strings.Split(date, "T")[0], close, high, low)
				}

				portfolios, signals := quantgo_algo.TurtleTradingAlgo(priceDatas, 5)
				xdates := dates[5:]
				xclosePrices := closePrices[5:]

				fmt.Printf("priceDatas: %v \n dates: %v \n portfolios: %v \n signals: %v \n xdates: %v",
					priceDatas, dates, portfolios, signals, xdates)

				fig := &grob.Fig{
					Data: []types.Trace{
						&grob.Scatter{
							X: types.DataArray(xdates),
							Y: types.DataArray(portfolios),
						},
					},
					Layout: &grob.Layout{
						Title: &grob.LayoutTitle{
							Text: types.StringType(fmt.Sprintf("A Figure of %s stock Portfolios with %s trading algorithm", quantStockcode, quantAlgo)),
						},
					},
				}

				offline.Show(fig)

				annotations := []grob.LayoutAnnotation{}

				for i, sig := range signals {
					if sig == 1 {
						anno := grob.LayoutAnnotation{}
						anno.X = xdates[i]
						anno.Y = xclosePrices[i]
						anno.Text = "buy in"
						annotations = append(annotations, anno)
					}
					if sig == -1 {
						anno := grob.LayoutAnnotation{}
						anno.X = xdates[i]
						anno.Y = xclosePrices[i]
						anno.Text = "buy out"
						annotations = append(annotations, anno)
					}
				}

				figSignals := &grob.Fig{
					Data: []types.Trace{
						&grob.Scatter{
							X: types.DataArray(dates),
							Y: types.DataArray(closePrices),
						},
					},

					Layout: &grob.Layout{
						Title: &grob.LayoutTitle{
							Text: types.StringType(fmt.Sprintf("A Figure of %s stock history", quantStockcode)),
						},
						Annotations: annotations,
					},
				}

				offline.Show(figSignals)
			}
		}
	},
}
