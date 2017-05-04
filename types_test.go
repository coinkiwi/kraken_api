package kraken

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_OHLCEntry_UnmarshalJSON(t *testing.T) {
	const incomingJSON = `[1493786460,"1326.860","1326.880","1324.533","1326.880","1326.643","3.93936569",9]`
	var data OHLCEntry
	if err := json.Unmarshal([]byte(incomingJSON), &data); err != nil {
		t.Error(err)
	}

	if data.Count != 9 {
		t.Errorf("data.Count expected: 9, got: %d", data.Count)
	}

	if !data.Timestamp.Equal(time.Unix(1493786460, 0)) {
		t.Errorf("data.Timestamp expected: 1493786460, got: %d", data.Timestamp.Unix())
	}
}

func Test_TradeData_UnmarshalJSON(t *testing.T) {
	const incomingJSON = `["1425.26000","8.96796823",1493926357.0243,"s","m","miscA"]`

	expectedTime := time.Unix(1493926357, 243*int64(100000))

	var data Trade
	if err := json.Unmarshal([]byte(incomingJSON), &data); err != nil {
		t.Error(err)
		t.Fail()
	}
	if data.Price != "1425.26000" {
		t.Errorf("data.Price expected: 1425.26000, got: %s", data.Price)
	}
	if data.Volume != "8.96796823" {
		t.Errorf("data.Volume expected: 8.96796823, got: %s", data.Volume)
	}

	if !data.Timestamp.Equal(expectedTime) {
		t.Errorf("The timestamps do not match. Expected %d, got: %d",
			expectedTime.UnixNano(),
			data.Timestamp.UnixNano())
	}

}
