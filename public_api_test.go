package kraken

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_Kraken_GetServerTime(t *testing.T) {
	var k Kraken
	k.Init()

	r, err := k.GetServerTime()
	if err != nil {
		t.Error(err)
	}

	if r.Unixtime < 1493752708 {
		t.Errorf("Unixtime might be wrong, %d", r.Unixtime)
	}

	if r.Rfc1123 == "" {
		t.Error("Rfc1123 string is empty")
	}
}

func Test_Kraken_GetAssetsInfo(t *testing.T) {
	var k Kraken
	k.Init()

	r, err := k.GetAssetsInfo()
	if err != nil {
		t.Error(err)
	}

	testCases := map[string]string{
		"XXBT": "XBT",
		"ZUSD": "USD",
		"ZEUR": "EUR",
	}

	for k, v := range testCases {
		if (*r)[k].Altname != v {
			t.Errorf("%s currency is missing", v)
		}
	}
}

func Test_Kraken_GetTradablePairs(t *testing.T) {
	var k Kraken
	k.Init()

	r, err := k.GetTradablePairs()
	if err != nil {
		t.Error(err)
	}
	type TwoStrings struct {
		Base, Quote string
	}
	testCases := map[string]TwoStrings{
		XETHXXBT: {"XETH", "XXBT"},
		XETCZEUR: {"XETC", "ZEUR"},
		XETCZUSD: {"XETC", "ZUSD"},
		XXBTZEUR: {"XXBT", "ZEUR"},
		XXBTZUSD: {"XXBT", "ZUSD"},
	}

	for k, v := range testCases {
		if (*r)[k].Base != v.Base || (*r)[k].Quote != v.Quote {
			t.Errorf("%s currency pair is missing", k)
		}
	}
}

func Test_Kraken_GetTickerInfo(t *testing.T) {
	var k Kraken
	k.Init()

	testCases := []string{XETHXXBT, XXBTZEUR}
	r, err := k.GetTickerInfo(testCases)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(*r) != 2 {
		t.Errorf("We should get back 2 TickerInfo structs. We got %d", len(*r))
	}
	if _, ok := (*r)[testCases[0]]; !ok {
		t.Errorf("%s return TickerInfo is missing", testCases[0])
	}
	if _, ok := (*r)[testCases[1]]; !ok {
		t.Errorf("%s return TickerInfo is missing", testCases[1])
	}
}

func Test_Kraken_GetOHLCData(t *testing.T) {
	var k Kraken
	k.Init()

	testCases := []string{"XETHXXBT", "XXBTZEUR"}
	for _, v := range testCases {
		op := NewOHLCQueryOptions()
		op.Pair = v
		r, err := k.GetOHLCData(op)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if r.Pair != v {
			t.Errorf("OHLC pair for %s is wrong", r.Pair)
		}
		if r.Last < 1493829480 {
			t.Errorf("We should get back a last OHLC data index. We got %d", r.Last)
		}
		if len(r.Data) < 2 {
			t.Errorf("OHLC data for %s is missing", r.Pair)
		}
	}
}

func Test_Kraken_GetOrderBook(t *testing.T) {
	var k Kraken
	k.Init()

	testCases := []string{"XETHXXBT", "XXBTZEUR"}
	for _, v := range testCases {
		obm, err := k.GetOrderBook(v, 10)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		ob, ok := (*obm)[v]
		if !ok {
			t.Errorf("Returned OrderBook doesn't have the %s pair data", v)
		}
		if ob.Pair != v {
			t.Errorf("Returned OrderBook pair mismatch. Expected %s, got: %s", v, ob.Pair)
		}
		if ob.Asks[0].Timestamp.UnixNano() < time.Unix(1493829480, 0).UnixNano() {
			t.Errorf("We should get back reasonably recent order. We got %v", ob.Asks[0].Timestamp)
		}
	}
}

func Test_Kraken_GetTrades(t *testing.T) {
	var k Kraken
	k.Init()

	const testCase = `{"error":[],"result":{"XXBTZEUR":[["1425.26000","8.96796823",1493926357.0243,"s","m",""],
	["1425.01000","0.01000000",1493926357.0391,"s","m",""],["1425.00000","0.10000000",1493926357.0579,"s","m",
	""],["1425.00000","0.10000000",1493926357.0624,"s","m",""]], "last":"1493926890306801911"}}`

	var trRaw TradesResult
	err := json.Unmarshal([]byte(testCase), &trRaw)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	tr := trRaw.Result
	if tr.Pair != "XXBTZEUR" {
		t.Errorf("Expected Pair should be XXBTZEUR, but got: %s", tr.Pair)
	}
	if tr.Last != "1493926890306801911" {
		t.Errorf("Expected Last should be 1493926890306801911, but got: %s", tr.Last)
	}
	if tr.Data[0].Volume != "8.96796823" {
		t.Errorf("Expected volume for first entry 8.96796823, but got: %s", tr.Data[0].Volume)
	}
	if tr.Data[3].Price != "1425.00000" {
		t.Errorf("Expected price for third entry 1425.00000, but got: %s", tr.Data[3].Price)
	}
	if tr.Data[3].Volume != "0.10000000" {
		t.Errorf("Expected volume for first entry 0.10000000, but got: %s", tr.Data[3].Volume)
	}
}
