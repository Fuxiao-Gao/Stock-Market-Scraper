package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"github.com/playwright-community/playwright-go"
)

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

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("Could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("Could not create page: %v", err)
	}

	for _, ticker := range tickers {
		url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", ticker)
		if _, err := page.Goto(url); err != nil {
			log.Printf("Could not navigate to %s: %v", url, err)
			continue
		}
		
		// --------Company Name Extraction with Wait--------
		companyLocator := page.Locator("h1.yf-4vbjci")

		// Wait for the element to be attached (or visible)
		if err := companyLocator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(5000),
		}); err != nil {
			log.Printf("Could not extract company name for %s: %v", ticker, err)
			continue
		}

		// Now safely get the text
		company, err := companyLocator.TextContent()
		if err != nil {
			log.Printf("Could not get company text for %s: %v", ticker, err)
			continue
		}

		// --------Price Extraction with Wait--------
		// Wait for the whole container that wraps price data
		// container := page.Locator("div.container.yf-16vvaki").First()
		// html, _ := container.InnerHTML()
		// log.Printf("Container HTML for %s: %s\n", ticker, html)
		// if err := container.WaitFor(playwright.LocatorWaitForOptions{
		// 	State:   playwright.WaitForSelectorStateVisible,
		// 	Timeout: playwright.Float(10000),
		// }); err != nil {
		// 	log.Printf("Price container not found for %s: %v", ticker, err)
		// 	continue
		// }

		priceLocator := page.Locator("fin-streamer[data-testid='qsp-price']").First()
		count, _ := priceLocator.Count()
		log.Printf("Found %d price elements for %s\n", count, ticker)
		if err := priceLocator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(5000),
		}); err != nil {
			log.Printf("Could not extract price for %s: %v", ticker, err)
			continue
		}

		price, err := priceLocator.TextContent()
		if err != nil {
			log.Printf("Could not get price text for %s: %v", ticker, err)
			continue
		}


		// --------Change Extraction with Wait--------
		changeLocator := page.Locator("fin-streamer[data-testid='qsp-price-change-percent']").First()
		count, _ = changeLocator.Count()
		log.Printf("Found %d change elements for %s\n", count, ticker)
		if err := changeLocator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(5000),
		}); err != nil {
			log.Printf("Could not extract price change percentage for %s: %v", ticker, err)
			continue
		}

		change, err := changeLocator.TextContent()
		if err != nil {
			log.Printf("Could not get price change text for %s: %v", ticker, err)
			continue
		}

		// Append the extracted data to the stocks slice
		stocks = append(stocks, Stock{
			Company: company,
			Price:   price,
			Change:  change,
		})
	}

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalf("could not create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Company", "Price", "Change"}
	if err := writer.Write(headers); err != nil {
		log.Fatalf("could not write headers to CSV: %v", err)
	}

	for _, stock := range stocks {
		record := []string{stock.Company, stock.Price, stock.Change}
		if err := writer.Write(record); err != nil {
			log.Printf("could not write record to CSV: %v", err)
		}
	}

	log.Println("CSV file written successfully!")
}
