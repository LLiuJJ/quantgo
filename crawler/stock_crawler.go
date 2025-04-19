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

package quantgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type DailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TimeSeriesDaily map[string]DailyData

type StockDataResponse struct {
	MetaData   MetaData        `json:"Meta Data"`
	TimeSeries TimeSeriesDaily `json:"Time Series (Daily)"`
}

func FetchStockDaliyHistoryFromRemote(code string, apiToken string) StockDataResponse {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s", code, apiToken)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var remoteDataResp StockDataResponse
	err = json.Unmarshal(body, &remoteDataResp)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	for date, data := range remoteDataResp.TimeSeries {
		fmt.Printf("  Date: %s\n", date)
		fmt.Printf("    Open: %s\n", data.Open)
		fmt.Printf("    High: %s\n", data.High)
		fmt.Printf("    Low: %s\n", data.Low)
		fmt.Printf("    Close: %s\n", data.Close)
		fmt.Printf("    Volume: %s\n", data.Volume)
	}
	return remoteDataResp
}
