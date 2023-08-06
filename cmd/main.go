package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"smartlab-dividend-fetcher/internal/cache"
	"smartlab-dividend-fetcher/internal/domain/dividend"
	"smartlab-dividend-fetcher/internal/domain/tickers"
	"strings"
)

const cacheFilename = ".dividends.cache"

func main() {
	const defaultTickerList = "ticker.list"
	var (
		printTickers   bool
		ticker         string
		tickerFile     string
		sortDescending bool
		noCache        bool
		verbose        bool

		outputStrategy = dividend.UpcomingOutputStrategy
		sortField      = dividend.TickerSortField
	)

	flag.BoolVar(&printTickers, "print-tickers", false, "Print hardcoded ticker map")
	flag.StringVar(&ticker, "ticker", "", "Fetch all history by specified ticker")
	flag.StringVar(&tickerFile, "f", defaultTickerList, "Path to file with ticker list, see ticker.list.example")
	flag.BoolVar(&noCache, "no-cache", false, "Don't read from cache")
	flag.BoolVar(&verbose, "v", false, "Verbose output")

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

	if !verbose {
		log.SetOutput(io.Discard)
	}

	var (
		err        error
		tickerList []string
	)
	if ticker != "" {
		tickerList = []string{strings.ToUpper(ticker)}
	} else {
		tickerList, err = tickers.ReadTickers(tickerFile)
		if err != nil {
			log.Fatalf("Can't read ticker file: %v\n", err)
		}
	}

	if printTickers {
		tickers.PrintTickers(os.Stdout, tickerList)
		return
	}

	fetcher := dividend.Fetcher{}
	cacheOpts := []cache.Option{}
	if noCache {
		cacheOpts = append(cacheOpts, cache.TTL(0))
	}
	fetcher.Cache, err = cache.NewCache(cacheFilename, cacheOpts...)
	if err != nil {
		log.Printf("Error opening cache, will do without it: %v", err)
		fetcher.Cache = &cache.DummyCache{}
	}

	dividends := fetcher.FetchDividends(tickerList...)
	if len(dividends) == 0 {
		fmt.Printf("No info on dividends")
		return
	}

	dividend.Output(os.Stdout, dividends, outputStrategy, sortField, sortDescending)
}
