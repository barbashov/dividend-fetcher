package main

import (
	"flag"
	"fmt"
	"os"
	"smartlab-dividend-fetcher/internal/domain/dividend"
	"smartlab-dividend-fetcher/internal/domain/tickers"
	"strings"
)

func main() {
	const defaultTickerList = "ticker.list"
	var (
		printTickers   bool
		ticker         string
		tickerFile     string
		sortDescending bool

		outputStrategy = dividend.UpcomingOutputStrategy
		sortField      = dividend.TickerSortField
	)

	flag.BoolVar(&printTickers, "print-tickers", false, "Print hardcoded ticker map")
	flag.StringVar(&ticker, "ticker", "", "Fetch all history by specified ticker")
	flag.StringVar(&tickerFile, "f", defaultTickerList, "Path to file with ticker list, see ticker.list.example")

	flag.Var(
		&outputStrategy,
		"out",
		fmt.Sprintf(
			"Output strategy. Posibble values: %s",
			strings.Join(dividend.AllOutputStrategies, ", "),
		),
	)

	flag.Var(
		&sortField,
		"sort",
		fmt.Sprintf(
			"Specifies sort field. Possible values: %s",
			strings.Join(dividend.AllSortFields, ", "),
		),
	)
	flag.BoolVar(&sortDescending, "desc", false, "Reverse sorting order")

	flag.Parse()

	var tickerList []string
	if ticker != "" {
		tickerList = []string{strings.ToUpper(ticker)}
	} else {
		var err error
		tickerList, err = tickers.ReadTickers(tickerFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Can't read ticker file: %v", err)
			os.Exit(1)
		}
	}

	if printTickers {
		tickers.PrintTickers(os.Stdout, tickerList)
		return
	}

	dividends := dividend.FetchDividends(tickerList...)
	if len(dividends) == 0 {
		fmt.Printf("No info on dividends")
		return
	}

	dividend.Output(os.Stdout, dividends, outputStrategy, sortField, sortDescending)
}
