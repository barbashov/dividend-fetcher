# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- Build: `make build` (produces `./main`) or `go build cmd/main.go`
- Tests: `make test` or `go test ./...`; single test: `go test ./internal/domain/dividend -run TestName`
- Lint: `make lint` (runs `golangci-lint` v1.53.3 via Docker)
- Run directly: `go run cmd/main.go [flags]`

## Architecture

Single-binary Go CLI that scrapes dividend tables from `smart-lab.ru` HTML pages (there is no API — parsing is literal string splitting on `<td>...</td>` markers in `internal/domain/dividend/fetcher.go`). Expect the scraper to break whenever the upstream HTML changes.

Flow: `cmd/main.go` → `tickers.ReadTickers` (whitespace-separated ticker file, default `ticker.list`) → `dividend.Fetcher.FetchDividends` → `dividend.Output`.

- `internal/domain/tickers` — ticker file reader and URL builder. `NormalizeTicker` strips the trailing `P` from a hardcoded list of preferred-share tickers (e.g. `SBERP` → `SBER`) because smart-lab.ru's URL scheme groups both share classes under the common-share page.
- `internal/domain/dividend` — fetcher, HTML parser, domain `Dividends` model (`model.go`), sort fields (`sort.go`), and output strategies (`output.go`: `upcoming`, `last`, `last-year`, `all`). The `OutputStrategy` and `SortField` types implement `flag.Value` so they validate CLI input in `main.go`.
- `internal/cache` — BoltDB-backed per-ticker cache at `.dividends.cache` with a 7-day default TTL. Cache miss / parse error / open failure all fall back gracefully (`DummyCache` or re-fetch). `-no-cache` forces TTL=0; `-v` enables the otherwise-discarded `log` output that traces cache hits/misses.

Adding a new output strategy requires: a constant in `output.go`, an entry in `AllOutputStrategies`, and a branch in `filterForOutput`. Adding a sort field follows the same pattern in `sort.go`.
