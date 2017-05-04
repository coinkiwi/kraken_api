package kraken

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// GetServerTime returns server time.
// Note: This is to aid in approximating the skew time between the server and client.
//
// https://www.kraken.com/help/api#get-server-time
func (k *Kraken) GetServerTime() (*ServerTime, error) {

	req, err := http.NewRequest("GET", urlGetServerTime, nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var dat ServerTimeResult
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	return &dat.Result, nil
}

// GetAssetsInfo returns assets info.
// https://www.kraken.com/help/api#get-asset-info
func (k *Kraken) GetAssetsInfo() (*AssetsInfoMap, error) {

	req, err := http.NewRequest("GET", urlGetAssetsInfo, nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var dat AssetsInfoResult
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	if len(dat.Error) > 0 {
		return nil, errors.New("JSON Error: " + dat.Error[0])
	}

	return &dat.Result, nil
}

// GetTradablePairs returns all tradable pairs from the api.
//
// Note: If an asset pair is on a maker/taker fee schedule, the taker side
// is given in "fees" and maker side in "fees_maker".
// For pairs not on maker/taker, they will only be given in "fees".
//
// https://www.kraken.com/help/api#get-tradable-pairs
func (k *Kraken) GetTradablePairs() (*AssetPairMap, error) {

	req, err := http.NewRequest("GET", urlGetTradablePairs, nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var dat AssetPairResult
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	if len(dat.Error) > 0 {
		return nil, errors.New("JSON Error: " + dat.Error[0])
	}

	return &dat.Result, nil
}

// GetTickerInfo return ticker info.
//
// Input: comma delimited list of asset pairs to get info on
// Result: array of pair names and their ticker info
//
// https://www.kraken.com/help/api#get-ticker-info
func (k *Kraken) GetTickerInfo(pairs []string) (*TickerInfoMap, error) {

	if len(pairs) == 0 {
		return nil, errors.New("JSON Error: Parameter pairs cannot be empty")
	}
	var pairList string
	pairList = pairs[0]
	for i := 1; i < len(pairs); i++ {
		pairList += "," + pairs[i]
	}

	req, err := http.NewRequest("GET", urlGetTickerInfo, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("pair", pairList)
	req.URL.RawQuery = query.Encode()

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var dat TickerInfoResult
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	if len(dat.Error) > 0 {
		return nil, errors.New("JSON Error: " + dat.Error[0])
	}

	return &dat.Result, nil
}

// GetOHLCData returns the OHLC data record for a given currency pair.
//
// Input:
// 	* pair = asset pair to get OHLC data for,
// 	* interval = time frame interval in minutes (optional): 1 (default), 5, 15, 30, 60, 240, 1440, 10080, 21600
// 	* since = return committed OHLC data since given id (optional.  exclusive)
//
// Output: array of pair name and OHLC data
// 	<pair_name> = pair name
//    	array of array entries(<time>, <open>, <high>, <low>, <close>, <vwap>, <volume>, <count>)
// 	last = id to be used as since when polling for new, committed OHLC data
//
// Note: the last entry in the OHLC array is for the current,
// not-yet-committed frame and will always be present, regardless of the value of "since".
//
// https://www.kraken.com/help/api#get-ohlc-data
func (k *Kraken) GetOHLCData(options *OHLCQueryOptions) (*OHLCEntryData, error) {

	if len(options.Pair) == 0 {
		return nil, errors.New("JSON Error: Parameter pair cannot be empty")
	}

	req, err := http.NewRequest("GET", urlGetOHLCData, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("pair", options.Pair)
	query.Add("interval", strconv.Itoa(options.Interval))
	if options.Since != "" {
		query.Add("since", options.Since)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// The OHLC data is not of uniform type, and needs to be processed manually
	// The format is a map of PAIR into the array of OHLCEntries, followed by "last": timestamp

	var dat map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}
	// check the error
	if errStringRaw, ok := dat["Error"]; ok {
		var errString []string
		if err := json.Unmarshal(errStringRaw, &errString); err != nil {
			return nil, err
		}
		if len(errString) > 0 {
			return nil, errors.New("JSON Error: " + errString[0])
		}
	}
	// parse the actual OHLC data entries
	ohlcDataMapRaw, ok := dat["result"]
	if !ok {
		return nil, errors.New("JSON Error: OHLC data result not present")
	}
	var ohlcDataMap map[string]json.RawMessage
	if err := json.Unmarshal(ohlcDataMapRaw, &ohlcDataMap); err != nil {
		return nil, err
	}
	ohlcData := &OHLCEntryData{}
	ohlcData.Pair = options.Pair
	tmpTimestamp, ok := ohlcDataMap["last"]
	if !ok {
		return nil, errors.New("JSON parsing error: missing 'last' property in OHLC data")
	}
	if err := json.Unmarshal(tmpTimestamp, &ohlcData.Last); err != nil {
		return nil, err
	}
	tmpData, ok := ohlcDataMap[options.Pair]
	if !ok {
		return nil, errors.New("JSON parsing error: missing 'pair' property in OHLC data")
	}
	if err := json.Unmarshal(tmpData, &ohlcData.Data); err != nil {
		return nil, err
	}

	return ohlcData, nil
}

// GetOrderBook returns current entries in the order book
//
// Input:
// 	pair = asset pair to get market depth for
//	count = maximum number of asks/bids
// Result:
//	<pair_name> = pair name
//	asks = ask side array of array entries(<price>, <volume>, <timestamp>)
//	bids = bid side array of array entries(<price>, <volume>, <timestamp>)
//
// https://www.kraken.com/help/api#get-order-book
func (k *Kraken) GetOrderBook(pair string, count int) (*OrderBookMap, error) {

	if len(pair) == 0 {
		return nil, errors.New("JSON Error: Parameter pair cannot be empty")
	}

	req, err := http.NewRequest("GET", urlGetOrderBook, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("pair", pair)
	query.Add("count", strconv.Itoa(count))
	req.URL.RawQuery = query.Encode()

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var dat OrderBookResult
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	if len(dat.Error) > 0 {
		return nil, errors.New("JSON Error: " + dat.Error[0])
	}

	tmp := dat.Result[pair]
	tmp.Pair = pair
	dat.Result[pair] = tmp

	return &dat.Result, nil
}
