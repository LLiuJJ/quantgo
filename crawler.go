package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

type FundDaliyData struct {
	DateTime           string
	FundCode           string
	FundName           string
	UnitNetValue       float64
	CumulativeNetValue float64
	DailyGrowthRate    float64
}

func FetchFundName(fundCode string) (string, error) {
	// 新浪财经单个基金详情页的URL模板（请根据实际情况调整）
	url := fmt.Sprintf("http://finance.sina.com.cn/fund/quotes/%s/bc.shtml", fundCode)

	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查状态码是否为200
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("请求失败: %s", resp.Status)
	}

	// 使用goquery加载文档
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 根据新浪财经的实际网页结构调整选择器
	var fundName string
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		// 假设基金名称位于<title>标签中，格式如：“基金名称_基金代码 - 数据中心-新浪财经”
		titleText := s.Text()
		parts := strings.Split(titleText, "_")
		if len(parts) > 0 {
			nameParts := strings.Split(parts[1], "-")
			if len(nameParts) > 0 {
				fundName = strings.TrimSpace(nameParts[0])
			}
		}
	})

	if fundName == "" {
		return "", fmt.Errorf("未找到基金名称")
	}

	return fundName, nil
}

func FetchFundData(fundCode string) []*FundDaliyData {

	fundDaliyData := []*FundDaliyData{}

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, -2, 0).Format("2006-01-02")

	fundName, _ := FetchFundName(fundCode)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("同步历史开始时间： %s 结束时间：\n", startDate, endDate)

	for i := 1; i < 20; i++ {
		url := fmt.Sprintf("http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=%s&page=%d&sdate=%s&edate=%s", fundCode, i, startDate, endDate)

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Referer", fmt.Sprintf("http://fund.eastmoney.com/%s.html", fundCode))

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		hasNext := true
		doc.Find("tbody tr").Each(func(index int, item *goquery.Selection) {
			tds := item.Find("td")
			date := strings.TrimSpace(tds.Eq(0).Text())
			unitNetValue := strings.TrimSpace(tds.Eq(1).Text())
			cumulativeNetValue := strings.TrimSpace(tds.Eq(2).Text())
			dailyGrowthRate := strings.TrimSpace(tds.Eq(3).Text())

			fmt.Printf("日期: %s, 单位净值: %s, 累计净值: %s, 日增长率: %s\n",
				date, unitNetValue, cumulativeNetValue, dailyGrowthRate)
			unitNetValueFloat, err := strconv.ParseFloat(unitNetValue, 64)
			if err != nil {
				fmt.Println("转换出错:", err)
			}
			cumulativeNetValueFloat, err := strconv.ParseFloat(cumulativeNetValue, 64)
			if err != nil {
				fmt.Println("转换出错:", err)
			}
			dailyGrowthRateFloat, err := strconv.ParseFloat(strings.TrimSuffix(dailyGrowthRate, "%"), 64)
			if err != nil {
				fmt.Println("转换出错:", err)
			}
			if strings.Contains(date, "暂无数据!") {
				hasNext = false
			} else {
				fundDaliyData = append(fundDaliyData, &FundDaliyData{FundCode: fundCode, FundName: fundName, DateTime: date, UnitNetValue: unitNetValueFloat, CumulativeNetValue: cumulativeNetValueFloat, DailyGrowthRate: dailyGrowthRateFloat})
			}
		})
		time.Sleep(time.Second * 2)
		if !hasNext {
			break
		}
	}

	return fundDaliyData
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
