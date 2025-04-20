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

package tradingalgo

import (
	"strconv"

	quantgo_crawler "github.com/LLiuJJ/quantgo/crawler"
)

// The Turtle Trading Algorithm is a trend-following strategy, and its core idea is:

// Buy Signal: Buy when the price breaks above the highest point of the previous N days.
// Sell Signal: Sell when the price falls below the lowest point of the previous N days.
// Stop-Loss Mechanism: Set a stop-loss point to control risk.

func TurtleTradingAlgo(data []quantgo_crawler.DailyData, n int) ([]float64, []int) {
	var portfolioValue []float64
	var buySellSignals []int
	cash := 100000.0
	shares := 0.0

	for i := n; i < len(data); i++ {
		highest := data[i-n].High
		lowest := data[i-n].Low
		for j := i - n + 1; j < i; j++ {
			if data[j].High > highest {
				highest = data[j].High
			}
			if data[j].Low < lowest {
				lowest = data[j].Low
			}
		}

		if data[i].Close > highest && shares == 0 {
			closeValue, err := strconv.ParseFloat(data[i].Close, 64)
			if err != nil {
				panic(err)
			}
			shares = cash / closeValue
			cash = 0
			buySellSignals = append(buySellSignals, 1)
		} else if data[i].Close < lowest && shares > 0 {
			closeValue, err := strconv.ParseFloat(data[i].Close, 64)
			if err != nil {
				panic(err)
			}
			cash = shares * closeValue
			shares = 0
			buySellSignals = append(buySellSignals, -1)
		} else {
			buySellSignals = append(buySellSignals, 0)
		}
		closeValue, err := strconv.ParseFloat(data[i].Close, 64)
		if err != nil {
			panic(err)
		}
		portfolioValue = append(portfolioValue, cash+shares*closeValue)
	}

	return portfolioValue, buySellSignals
}
