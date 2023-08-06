package tickers

import (
	"io"
	"os"
	"regexp"

	"github.com/olekukonko/tablewriter"
)

var splitRe = regexp.MustCompile("[\n\r\t ]+")

func ReadTickers(filename string) ([]string, error) {
	d, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return splitRe.Split(string(d), -1), nil
}

func PrintTickers(w io.Writer, tickers []string) {
	t := tablewriter.NewWriter(w)
	t.SetHeader([]string{"Ticker", "Normalized", "URL"})
	t.SetBorder(false)
	for _, ticker := range tickers {
		t.Append([]string{
			ticker,
			NormalizeTicker(ticker),
			BuildFetchLink(ticker),
		})
	}
	t.Render()
}
