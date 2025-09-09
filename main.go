package main

import (
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"github.com/gocolly/colly/v2"
)

type Stock struct {
	company, price, change string
}

func main() {
	ticker := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAP",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := []Stock{}
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Look for <div id="quote-header-info"> ... </div> in the page.
	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}
		// Inside the div#quote-header-info, find the <h1> element.
		stock.company = e.ChildText("h1")
		fmt.Println("Company", stock.company)
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("Price", stock.price)
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		fmt.Println("Change", stock.change)

		stocks = append(stocks, stock)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response length:", len(r.Body))
		fmt.Println(string(r.Body[:1000])) // first 1000 chars
	})


	c.Wait()

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}


	fmt.Println(stocks)
	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("failed to create the output csv file", err)
	}

	writer := csv.NewWriter(file)
	// write the headers to the csv file
	headers := []string{
		"company",
		"price",
		"change",
	}
	writer.Write(headers)

	// write each stock as a record in the csv file
	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}

	// in go, defers run in the LIFO order when written separately
	// in this case, defer schedules the entire function to run at the end of the current function, not immediately
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}

		if err := file.Close(); err != nil {
			log.Fatalln("error closing file:", err)
		}
		fmt.Println("done")
	}()

}