package kraken

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	urlBaseURL          string = "https://api.kraken.com/0/public/"
	urlGetServerTime    string = urlBaseURL + "Time"
	urlGetAssetsInfo    string = urlBaseURL + "Assets"
	urlGetTradablePairs string = urlBaseURL + "AssetPairs"
	urlGetTickerInfo    string = urlBaseURL + "Ticker"
	urlGetOHLCData      string = urlBaseURL + "OHLC"
	urlGetOrderBook     string = urlBaseURL + "Depth"
	urlGetTrades        string = urlBaseURL + "Trades"
)

/* Some of the common pairs for convenience. */
const (
	XETHXXBT = "XETHXXBT"
	XETHZCAD = "XETHZCAD"
	XETHZEUR = "XETHZEUR"
	XETHZGBP = "XETHZGBP"
	XETHZUSD = "XETHZUSD"

	XETCXXBT = "XETCXXBT"
	XETCZCAD = "XETCZCAD"
	XETCZEUR = "XETCZEUR"
	XETCZGBP = "XETCZGBP"
	XETCZUSD = "XETCZUSD"

	XLTCZCAD = "XLTCZCAD"
	XLTCZEUR = "XLTCZEUR"
	XLTCZUSD = "XLTCZUSD"
	XXBTXLTC = "XXBTXLTC"
	XXBTZCAD = "XXBTZCAD"
	XXBTZEUR = "XXBTZEUR"
	XXBTZGBP = "XXBTZGBP"
	XXBTZUSD = "XXBTZUSD"
)

// APIError represents JSON error type.
type APIError []string

// ServerTime contains server unix timestamp and string date.
type ServerTime struct {
	Unixtime int64  `json:"unixtime"`
	Rfc1123  string `json:"rfc1123"`
}

// Human readable printout for the serverTime
func (t ServerTime) String() string {
	return fmt.Sprintf("unixtime: %d,  rfc1123: %s",
		t.Unixtime, t.Rfc1123)
}

// ServerTimeResult result from the JSON call.
type ServerTimeResult struct {
	Result ServerTime `json:"result"`
	Error  APIError   `json:"error"`
}

// AssetInfo contains details about currency.
type AssetInfo struct {
	Altname         string `json:"altname"`
	Aclass          string `json:"aclass"`
	Decimals        byte   `json:"decimals"`
	DisplayDecimals byte   `json:"display_decimals"`
}

// AssetsInfoMap maps string -> AssetInfo
type AssetsInfoMap map[string]AssetInfo

// Human readable printout for the AssetInfo
func (t AssetInfo) String() string {
	return fmt.Sprintf("altname: %s, aclass: %s, decimals: %d, display: %d",
		t.Altname, t.Aclass, t.Decimals, t.DisplayDecimals)
}

// AssetsInfoResult represents the result from the JSON call.
type AssetsInfoResult struct {
	Result AssetsInfoMap `json:"result"`
	Error  APIError      `json:"error"`
}

// FeeInfo is a tuple for fee info.
type FeeInfo []float32

// AssetPairInfo contains details about currency pair.
type AssetPairInfo struct {
	// Alternate pair name.
	Altname string `json:"altname"`
	// Asset class of base component.
	AclassBase string `json:"aclass_base"`
	// Asset id of base component.
	Base string `json:"base"`
	// Asset class of quote component.
	AclassQuote string `json:"aclass_quote"`
	// Asset id of quote component.
	Quote string `json:"quote"`
	// Volume lot size.
	Lot string `json:"lot"`
	// Scaling decimal places for pair.
	PairDecimals byte `json:"pair_decimals"`
	// Scaling decimal places for volume.
	LotDecimals byte `json:"lot_decimals"`
	// Amount to multiply lot volume by to get currency volume.
	LotMultiplier byte `json:"lot_multiplier"`
	// Array of leverage amounts available when buying.
	LeverageBuy []byte `json:"leverage_buy"`
	// Array of leverage amounts available when selling.
	LeverageSell []byte `json:"leverage_sell"`
	// Fee schedule array in [volume, percent fee] tuples.
	Fees []FeeInfo `json:"fees"`
	// Maker fee schedule array in [volume, percent fee] tuples (if on maker/taker).
	FeesMaker []FeeInfo `json:"fees_maker"`
	// Volume discount currency.
	FeeVolumeCurrency string `json:"fee_volume_currency"`
	// Margin call level.
	MarginCall byte `json:"margin_call"`
	// Stop-out/liquidation margin level.
	MarginStop byte `json:"margin_stop"`
}

// AssetPairMap maps AssetsPair data to currency pair.
type AssetPairMap map[string]AssetPairInfo

// AssetPairResult represents the result from the JSON API call.
type AssetPairResult struct {
	Result AssetPairMap `json:"result"`
	Error  APIError     `json:"error"`
}

// TickerInfo contains the current ticker info.
type TickerInfo struct {
	// Ask array(<price>, <whole lot volume>, <lot volume>).
	A []string
	// Bid array(<price>, <whole lot volume>, <lot volume>).
	B []string
	// Last trade closed array(<price>, <lot volume>).
	C []string
	// Volume array(<today>, <last 24 hours>).
	V []string
	// Volume weighted average price array(<today>, <last 24 hours>).
	P []string
	// Number of trades array(<today>, <last 24 hours>).
	T []int
	// Low array(<today>, <last 24 hours>).
	L []string
	// High array(<today>, <last 24 hours>).
	H []string
	// Today's opening price.
	O string
}

// TickerInfoMap maps currency to TickerInfo
type TickerInfoMap map[string]TickerInfo

// TickerInfoResult result from the JSON API call.
type TickerInfoResult struct {
	Result TickerInfoMap `json:"result"`
	Error  APIError      `json:"error"`
}

// OHLCEntry has a single OHLC entry.
type OHLCEntry struct {
	Timestamp time.Time
	Data      [6]string
	Count     int64
}

// UnmarshalJSON for custom marchaling of OHLCEntry
func (c *OHLCEntry) UnmarshalJSON(b []byte) error {
	tmp := [8]json.Number{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	tmpTimeSec, err := tmp[0].Int64()
	if err != nil {
		return err
	}
	c.Timestamp = time.Unix(tmpTimeSec, 0)

	var data [6]string
	for i := 1; i < 7; i++ {
		data[0] = tmp[i].String()
	}

	c.Data = data
	c.Count, err = tmp[7].Int64()
	if err != nil {
		return err
	}

	return nil
}

// OHLCEntryData contains the result from the JSON API call.
type OHLCEntryData struct {
	Data []OHLCEntry
	Pair string
	Last int64
}

// OHLCQueryOptions contains the query parameters for the OHLC data request.
type OHLCQueryOptions struct {
	Pair     string
	Interval int
	Since    string
}

// NewOHLCQueryOptions creates a new, default instance of OHLC request options.
//
// 	Default values:
//	* pair: XXBTZEUR
//	* interval: 1 minute
//	* since: <empty>
func NewOHLCQueryOptions() *OHLCQueryOptions {
	op := &OHLCQueryOptions{}
	op.Interval = 1
	op.Pair = XXBTZEUR
	op.Since = ""
	return op
}

// OrderBookEntry represents a single entry: price, volume, timestamp
type OrderBookEntry struct {
	Timestamp time.Time
	Price     string
	Volume    string
}

// UnmarshalJSON of the OrderBookEntry
func (o *OrderBookEntry) UnmarshalJSON(b []byte) error {
	tmp := [3]json.Number{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	tmpTime, err1 := tmp[2].Int64()
	if err1 != nil {
		return err1
	}
	o.Timestamp = time.Unix(tmpTime, 0)
	o.Price = tmp[0].String()
	o.Volume = tmp[1].String()
	return nil
}

// OrderBook the bids and asks for a given currency pair
type OrderBook struct {
	Pair string
	Asks []OrderBookEntry
	Bids []OrderBookEntry
}

// UnmarshalJSON of the OrderBook
func (o *OrderBook) UnmarshalJSON(b []byte) error {
	var tmp = map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	nAsks := len(tmp["asks"])
	var asks = make([]OrderBookEntry, nAsks)
	if err := json.Unmarshal(tmp["asks"], &asks); err != nil {
		return err
	}
	o.Asks = asks
	nBids := len(tmp["bids"])
	var bids = make([]OrderBookEntry, nBids)
	if err := json.Unmarshal(tmp["bids"], &bids); err != nil {
		return err
	}
	o.Bids = bids[:]

	return nil
}

// OrderBookMap maps the currency pair to OrderBook bids and asks struct
type OrderBookMap map[string]OrderBook

// OrderBookResult result from the JSON API call.
type OrderBookResult struct {
	Result OrderBookMap `json:"result"`
	Error  APIError     `json:"error"`
}

// Trade represents single trade
type Trade struct {
	Timestamp time.Time
	Price     string
	Volume    string
	BS        string
	ML        string
	MISC      string
}

// UnmarshalJSON for the TradeData
func (t *Trade) UnmarshalJSON(b []byte) error {
	var tmp = [6]json.RawMessage{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &t.Price); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &t.Volume); err != nil {
		return err
	}
	var tmpTimeJSON json.Number
	if err := json.Unmarshal(tmp[2], &tmpTimeJSON); err != nil {
		return err
	}
	tmpTimeFloat, err := tmpTimeJSON.Float64()
	if err != nil {
		return err
	}
	var tmpTimeMilliseconds = int64(tmpTimeFloat * 10000)
	t.Timestamp = time.Unix(0, tmpTimeMilliseconds*int64(100000))
	if err := json.Unmarshal(tmp[3], &t.BS); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[4], &t.ML); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[5], &t.MISC); err != nil {
		return err
	}
	return nil
}

// TradeBook the bids and asks for a given currency pair
type TradeBook struct {
	Pair string
	Data []Trade
	Last string
}

// UnmarshalJSON of the TradeBook
func (o *TradeBook) UnmarshalJSON(b []byte) error {
	var tmp = map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for k, v := range tmp {
		if k == "last" {
			if err := json.Unmarshal(v, &o.Last); err != nil {
				return err
			}
		} else {
			// k must be the Pair
			o.Pair = k
			if err := json.Unmarshal(v, &o.Data); err != nil {
				return err
			}
		}
	}

	return nil
}

// TradesResult result from the JSON API call.
type TradesResult struct {
	Result TradeBook `json:"result"`
	Error  APIError  `json:"error"`
}
