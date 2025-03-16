package dividend

import (
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"time"

	"github.com/olekukonko/tablewriter"
)

type OutputStrategy string

const (
	UpcomingOutputStrategy OutputStrategy = "upcoming"
	LastOutputStrategy     OutputStrategy = "last"
	LastYearOutputStrategy OutputStrategy = "last-year"
	AllOutputStrategy      OutputStrategy = "all"
)

var AllOutputStrategies = []string{
	string(UpcomingOutputStrategy),
	string(LastOutputStrategy),
	string(LastYearOutputStrategy),
	string(AllOutputStrategy),
}

func Output(w io.Writer, divs map[string][]Dividends, outStrat OutputStrategy, sort SortField, sortDescending bool) {
	t := tablewriter.NewWriter(w)
	t.SetHeader([]string{
		"Ticker",
		"T+2",
		"ExDividend",
		"Period",
		"Dividend",
		"Price",
		"Yield",
		"Error",
	})
	t.SetBorder(false)

	lines := []Dividends{}
	for _, div := range divs {
		lines = append(lines, filterForOutput(div, outStrat)...)
	}

	if len(lines) == 0 {
		t.Append([]string{"No data"})
		t.Render()
		return
	}

	slices.SortFunc(lines, SortFunc(sort, sortDescending))

	for _, d := range lines {
		t.Append(prepareLine(d))
	}

	t.Render()
}

func filterForOutput(div []Dividends, outStrat OutputStrategy) []Dividends {
	if outStrat == LastOutputStrategy {
		return []Dividends{div[0]}
	}

	if outStrat == UpcomingOutputStrategy && div[0].ExDividend.After(time.Now()) {
		return []Dividends{div[0]}
	}

	if outStrat == AllOutputStrategy {
		return div
	}

	if outStrat == LastYearOutputStrategy {
		ret := []Dividends{}
		lastYear := time.Now().Add(-366 * 24 * time.Hour)
		for _, div := range div {
			if div.ExDividend.Before(lastYear) {
				continue
			}
			ret = append(ret, div)
		}
		return ret
	}

	return nil
}

func prepareLine(d Dividends) []string {
	serr := ""
	if d.Err != nil {
		serr = fmt.Sprintf("%v", d.Err)
	}
	return []string{
		d.Ticker,
		d.T_2.Format(TimeLayout),
		d.ExDividend.Format(TimeLayout),
		d.Period,
		fmt.Sprintf("%.04f", d.Dividend),
		fmt.Sprintf("%.02f", d.Price),
		fmt.Sprintf("%.02f%%", d.Yield),
		serr,
	}
}

func (o *OutputStrategy) String() string {
	return string(*o)
}

func (o *OutputStrategy) Set(new string) error {
	if !slices.Contains(AllOutputStrategies, new) {
		return fmt.Errorf("unknown output strategy %q", new)
	}
	*o = OutputStrategy(new)
	return nil
}
