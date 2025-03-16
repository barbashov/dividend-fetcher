package dividend

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"smartlab-dividend-fetcher/internal/domain/tickers"
	"strconv"
	"strings"
	"time"
)

type Fetcher struct {
	Cache Cache
}

func (f *Fetcher) FetchDividends(tickers ...string) map[string][]Dividends {
	ret := map[string][]Dividends{}
	for _, ticker := range tickers {
		if ticker == "" {
			continue
		}
		divs, err := f.Cache.Get(ticker)
		if err != nil {
			log.Printf("No info in %s in cache, will fetch", ticker)
			divs = f.fetchTicker(ticker)
			err = f.Cache.Set(ticker, divs)
			if err != nil {
				log.Printf("Can't save %s to cache: %v", ticker, err)
			} else {
				log.Printf("Saved %s to cache", ticker)
			}
		} else {
			log.Printf("Got %s info from cache", ticker)
		}
		if len(divs) == 0 {
			continue
		}
		ret[ticker] = divs
	}
	return ret
}

func (f *Fetcher) fetchTicker(ticker string) []Dividends {
	cl := http.Client{}
	url := tickers.BuildFetchLink(ticker)
	res, err := cl.Get(url)
	if err != nil {
		return []Dividends{{
			Ticker: ticker,
			Err:    err,
		}}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []Dividends{{
			Ticker: ticker,
			Err:    fmt.Errorf("invalid HTTP status %d fetching %s", res.StatusCode, url),
		}}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []Dividends{{
			Ticker: ticker,
			Err:    err,
		}}
	}

	s := string(data)
	els := strings.Split(s, fmt.Sprintf("<td>%s</td>", ticker))
	if len(els) < 2 { // no dividends
		return nil
	}

	cleanReplacer := strings.NewReplacer(
		"<td>", "",
		"</td>", "",
		"<strong>", "",
		"</strong>", "",
		"<strong >", "",
		"₽", "",
		"%", "",
		",", ".",
	)

	els = els[1:]
	ret := []Dividends{}
	for ri, el := range els {
		cols := strings.Split(el, "</td>")
		if len(cols) < 6 {
			ret = append(ret, Dividends{
				Err:    fmt.Errorf("row #%d: not enought columns", ri),
				Ticker: ticker,
			})
			continue
		}
		for i, col := range cols {
			cols[i] = cleanReplacer.Replace(col)
			cols[i] = strings.TrimSpace(cols[i])
		}

		if strings.Contains(cols[3], "dividend_canceled") {
			continue
		}

		errs := []error{}
		divs := Dividends{
			Ticker: ticker,
			Period: cols[2],
		}

		divs.T_2, err = time.Parse(TimeLayout, cols[0])
		if err != nil {
			errs = append(errs, fmt.Errorf("row #%d: error parsing time %s", ri, cols[0]))
		}
		divs.ExDividend, err = time.Parse(TimeLayout, cols[1])
		if err != nil {
			errs = append(errs, fmt.Errorf("row #%d: error parsing time %s", ri, cols[1]))
		}

		divs.Dividend, err = strconv.ParseFloat(cols[3], 64)
		if err != nil {
			errs = append(errs, fmt.Errorf("row #%d: error parsing number %s", ri, cols[3]))
		}
		divs.Price, err = strconv.ParseFloat(cols[4], 64)
		if err != nil {
			errs = append(errs, fmt.Errorf("row #%d: error parsing number %s", ri, cols[4]))
		}
		divs.Yield, err = strconv.ParseFloat(cols[5], 64)
		if err != nil {
			errs = append(errs, fmt.Errorf("row #%d: error parsing number %s", ri, cols[5]))
		}

		divs.Err = errors.Join(errs...)
		ret = append(ret, divs)
	}

	return ret
}
