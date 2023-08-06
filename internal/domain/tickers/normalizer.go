package tickers

import (
	"fmt"
	"golang.org/x/exp/slices"
)

var tickersPrefs = []string{
	"BANEP", "BISVP", "BSPBP", "CNTLP", "DZRDP", "GAZAP", "HIMCP",
	"IGSTP", "JNOSP", "KAZTP", "KCHEP", "KGKCP", "KRKNP", "KRKOP",
	"KROTP", "KRSBP", "KTSBP", "KZOSP", "LNZLP", "LSNGP", "MAGEP",
	"MFGSP", "MGTSP", "MISBP", "MTLRP", "NKNCP", "NNSBP", "OMZZP",
	"PMSBP", "RTKMP", "RTSBP", "SAGOP", "SAREP", "SBERP", "SNGSP",
	"STSBP", "TASBP", "TATNP", "TGKBP", "TORSP", "TRNFP", "VGSBP",
	"VJGZP", "VRSBP", "VSYDP", "WTCMP", "YKENP", "YRSBP",
}

func NormalizeTicker(ticker string) string {
	if slices.Contains(tickersPrefs, ticker) {
		ticker = ticker[:4]
	}
	return ticker
}

func BuildFetchLink(ticker string) string {
	const urlPattern = "https://smart-lab.ru/q/%s/dividend/"
	return fmt.Sprintf(urlPattern, NormalizeTicker(ticker))
}
