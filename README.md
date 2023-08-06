# Dividend fetcher

Fetches dividend info from [smart-lab.ru](https://smart-lab.ru).

## Usage

File `internal/tickers/default.go` contains a list of default tickers to fetch.

Ran without any arguments, tool will fetch all tickers from the file and output upcoming dividend info.

Behaviour can be changed with flags:

```
  -desc
        Reverse sorting order
  -f string
        Path to file with ticker list, see ticker.list.example (default "ticker.list")
  -out value
        Output strategy. Posibble values: upcoming, last, all (default upcoming)
  -print-tickers
        Print hardcoded ticker map
  -sort value
        Specifies sort field. Possible values: ticker, t2, ex_dividend, period, dividend, price, yield (default ticker)
  -ticker string
        Fetch all history by specified ticker
```

### Example

```
% go run cmd/main.go -ticker GAZP -out all -sort yield -desc
  TICKER |    T+2     | EXDIVIDEND |  PERIOD  | DIVIDEND | PRICE  | YIELD   
---------+------------+------------+----------+----------+--------+---------
  GAZP   | 18.07.2022 | 20.07.2022 | 2021 год |   0.0000 | 187.15 | 28.10%  
  GAZP   | 07.10.2022 | 11.10.2022 | 2кв 2022 |  51.0300 | 195.00 | 26.20%  
  GAZP   | 14.07.2020 | 16.07.2020 | 2019 год |  15.2400 | 194.56 | 7.80%   
  GAZP   | 16.07.2019 | 18.07.2019 | 2018 год |  16.6100 | 238.01 | 7.00%   
  GAZP   | 18.07.2017 | 20.07.2017 | 2016 год |   8.0397 | 124.49 | 6.50%   
  GAZP   | 17.07.2018 | 19.07.2018 | 2017 год |   8.0400 | 145.31 | 5.50%   
  GAZP   | 20.07.2016 | 20.07.2016 | 2015 год |   7.8900 | 145.72 | 5.40%   
  GAZP   | 16.07.2015 | 16.07.2015 | 2014 год |   7.2000 | 146.18 | 4.90%   
  GAZP   | 17.07.2014 | 17.07.2014 | 2013 год |   7.2000 | 147.00 | 4.90%   
  GAZP   | 13.05.2013 | 13.05.2013 | 2012 год |   5.9900 | 129.70 | 4.60%   
  GAZP   | 13.07.2021 | 15.07.2021 | 2020 год |  12.5500 | 294.86 | 4.30%   
```