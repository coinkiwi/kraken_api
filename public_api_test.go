package kraken_api

import (
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
		"XETHXXBT": {"XETH", "XXBT"},
		"XETCZEUR": {"XETC", "ZEUR"},
		"XETCZUSD": {"XETC", "ZUSD"},
		"XXBTZEUR": {"XXBT", "ZEUR"},
		"XXBTZUSD": {"XXBT", "ZUSD"},
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

	testCases := []string{"XETHXXBT", "XXBTZEUR"}
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
