package kraken

import (
	"encoding/json"
	"testing"
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

func Test_OHLCEntry_UnmarshalJSON(t *testing.T) {
	const incomingJSON = `[1493786460,"1326.860","1326.880","1324.533","1326.880","1326.643","3.93936569",9]`
	var data OHLCEntry
	if err := json.Unmarshal([]byte(incomingJSON), &data); err != nil {
		t.Errorf("Got error when unmarshal was called: %s", err)
	}

	if data.Count != 9 {
		t.Errorf("data.Count expected: 9, got: %d", data.Count)
	}

	if data.Timestamp != 1493786460 {
		t.Errorf("data.Timestamp expected: 1493786460, got: %d", data.Timestamp)
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
		if ob.Asks[0].Timestamp < 1493829480 {
			t.Errorf("We should get back reasonably recent order. We got %d", ob.Asks[0].Timestamp)
		}
	}
}
