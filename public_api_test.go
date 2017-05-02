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
