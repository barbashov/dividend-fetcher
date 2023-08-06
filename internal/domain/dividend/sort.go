package dividend

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

type SortField string

const (
	TickerSortField     SortField = "ticker"
	T_2SortField        SortField = "t2"
	ExDividendSortField SortField = "ex_dividend"
	PeriodSortField     SortField = "period"
	DividendSortField   SortField = "dividend"
	PriceSortField      SortField = "price"
	YieldSortField      SortField = "yield"
)

var AllSortFields = []string{
	string(TickerSortField),
	string(T_2SortField),
	string(ExDividendSortField),
	string(PeriodSortField),
	string(DividendSortField),
	string(PriceSortField),
	string(YieldSortField),
}

func SortFunc(field SortField, descending bool) func(a, b Dividends) int {
	orderModified := 1
	if descending {
		orderModified = -1
	}

	sign := func(x float64) int {
		if x < 0 {
			return -1
		} else if x > 0 {
			return 1
		}
		return 0
	}

	return func(a, b Dividends) int {
		switch field {
		case TickerSortField:
			return strings.Compare(a.Ticker, b.Ticker) * orderModified
		case T_2SortField:
			return a.T_2.Compare(b.T_2) * orderModified
		case ExDividendSortField:
			return a.ExDividend.Compare(b.ExDividend) * orderModified
		case PeriodSortField:
			return strings.Compare(a.Period, b.Period) * orderModified
		case DividendSortField:
			return sign(a.Dividend-b.Dividend) * orderModified
		case PriceSortField:
			return sign(a.Price-b.Price) * orderModified
		case YieldSortField:
			return sign(a.Yield-b.Yield) * orderModified
		}
		return 0 // should panic?
	}
}

func (s *SortField) String() string {
	return string(*s)
}

func (s *SortField) Set(new string) error {
	if !slices.Contains(AllSortFields, new) {
		return fmt.Errorf("unknown sort field %q", new)
	}

	*s = SortField(new)
	return nil
}
