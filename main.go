package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type QuoteResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol                     string  `json:"symbol"`
			ShortName                  string  `json:"shortName"`
			RegularMarketPrice         float64 `json:"regularMarketPrice"`
			RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
		} `json:"result"`
	} `json:"quoteResponse"`
}

type Stock struct {
	Company string
	Price   string
	Change  string
}

func main() {
	tickers := []string{
		"MSFT", "IBM", "GE", "UNP", "COST",
		"MCD", "V", "WMT", "DIS", "MMM",
		"INTC", "AXP", "AAP", "BA", "CSCO",
		"GS", "JPM", "CRM", "VZ",
	}

	stocks := []Stock{}

	// Build the Yahoo Finance API URL
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s",
		strings.Join(tickers, ","),
	)

	// Make HTTP GET request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Failed to create request:", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
	"(KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Failed to fetch data:", err)
	}
	defer resp.Body.Close()

	// check HTTP status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected HTTP status: %s", resp.Status)
	}

	// Decode JSON
	var data QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalln("Failed to decode JSON:", err)
	}

	// Extract stock data
	for _, q := range data.QuoteResponse.Result {
		stock := Stock{
			Company: q.ShortName,
			Price:   fmt.Sprintf("%.2f", q.RegularMarketPrice),
			Change:  fmt.Sprintf("%.2f%%", q.RegularMarketChangePercent),
		}
		fmt.Println(stock.Company, stock.Price, stock.Change)
		stocks = append(stocks, stock)
	}

	// Write CSV
	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create CSV file:", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"company", "price", "change"}
	if err := writer.Write(headers); err != nil {
		log.Fatalln("Error writing headers:", err)
	}

	for _, s := range stocks {
		record := []string{s.Company, s.Price, s.Change}
		if err := writer.Write(record); err != nil {
			log.Println("Error writing record:", err)
		}
	}

	fmt.Println("CSV file written successfully!")
}
