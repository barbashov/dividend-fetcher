package dividend

import "time"

const TimeLayout = "02.01.2006"

type Dividends struct {
	Ticker     string
	T_2        time.Time
	ExDividend time.Time
	Period     string
	Dividend   float64
	Price      float64
	Yield      float64
	Err        error
}
