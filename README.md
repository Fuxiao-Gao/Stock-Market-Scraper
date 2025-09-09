# Stock Market Scraper
This project is a Go-based scraper that collects stock information(company name, current price, and daily percentage change) from Yahoo Finance and writes the results to a csv file.

# Features
- Fetch live stock data from Yahoo Finance
- Extract company name, price, and percentage change
- Write results to a structured csv file (stocks.csv)
- Handle multiple stock tickers in a single run
![Stock price screenshot](examples/example_csv.png)

# Requirements
- Go 1.19+
- Playwright for Go
Playwrigt is used to handle dynamic content on Yahoo Finance that cannnot be scraped reliably via static HTTP requests

# Currently working on adding:
- Schedule daily emails about stock price change
