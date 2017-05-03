package kraken_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const URLGetServerTime string = "https://api.kraken.com/0/public/Time"
const URLGetAssetsInfo string = "https://api.kraken.com/0/public/Assets"
const URLGetTradablePairs string = "https://api.kraken.com/0/public/AssetPairs"
const URLGetTickerInfo string = "https://api.kraken.com/0/public/Ticker"

/*
 Represents JSON error type.
*/
type KrakenApiError []string

/*
 Server time.
*/
type ServerTime struct {
	Unixtime uint64 `json:"unixtime"`
	Rfc1123  string `json:"rfc1123"`
}

func (t ServerTime) String() string {
	return fmt.Sprintf("unixtime: %d,  rfc1123: %s",
		t.Unixtime, t.Rfc1123)
}

type ServerTimeResult struct {
	Result ServerTime     `json:"result"`
	Error  KrakenApiError `json:"error"`
}

/*
 Get server time.
 Note: This is to aid in approximating the skew time between the server and client.

 https://www.kraken.com/help/api#get-server-time
*/
func (k *Kraken) GetServerTime() (*ServerTime, error) {

	req, err := http.NewRequest("GET", URLGetServerTime, nil)
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

/*
 Assets info.
 https://www.kraken.com/help/api#get-asset-info
*/
type AssetInfo struct {
	Altname         string `json:"altname"`
	Aclass          string `json:"aclass"`
	Decimals        byte   `json:"decimals"`
	DisplayDecimals byte   `json:"display_decimals"`
}

type AssetsInfoMap map[string]AssetInfo

func (t AssetInfo) String() string {
	return fmt.Sprintf("altname: %d, aclass: %d, decimals: %d, display: %d",
		t.Altname, t.Aclass, t.Decimals, t.DisplayDecimals)
}

type AssetsInfoResult struct {
	Result AssetsInfoMap  `json:"result"`
	Error  KrakenApiError `json:"error"`
}

/*
 Get assets info.

 https://www.kraken.com/help/api#get-asset-info
*/
func (k *Kraken) GetAssetsInfo() (*AssetsInfoMap, error) {

	req, err := http.NewRequest("GET", URLGetAssetsInfo, nil)
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
		return nil, errors.New("We got error" + dat.Error[0])
	}

	return &dat.Result, nil
}

type FeeInfo []float32

type AssetPairInfo struct {
	/* Alternate pair name. */
	Altname string `json:"altname"`
	/* Asset class of base component */
	AclassBase string `json:"aclass_base"`
	/* Asset id of base component */
	Base string `json:"base"`
	/* Asset class of quote component. */
	AclassQuote string `json:"aclass_quote"`
	/* Asset id of quote component. */
	Quote string `json:"quote"`
	/* Volume lot size. */
	Lot string `json:"lot"`
	/* Scaling decimal places for pair. */
	PairDecimals byte `json:"pair_decimals"`
	/* Scaling decimal places for volume. */
	LotDecimals byte `json:"lot_decimals"`
	/* Amount to multiply lot volume by to get currency volume. */
	LotMultiplier byte `json:"lot_multiplier";"`
	/* Array of leverage amounts available when buying. */
	LeverageBuy []byte `json:"leverage_buy"`
	/* Array of leverage amounts available when selling. */
	LeverageSell []byte `json:"leverage_sell"`
	/* Fee schedule array in [volume, percent fee] tuples. */
	Fees []FeeInfo `json:"fees"`
	/* Maker fee schedule array in [volume, percent fee] tuples (if on maker/taker). */
	FeesMaker []FeeInfo `json:"fees_maker"`
	/* Volume discount currency. */
	FeeVolumeCurrency string `json:"fee_volume_currency"`
	/* Margin call level. */
	MarginCall byte `json:"margin_call"`
	/* Stop-out/liquidation margin level. */
	MarginStop byte `json:"margin_stop"`
}

type AssetPairMap map[string]AssetPairInfo

type AssetPairResult struct {
	Result AssetPairMap   `json:"result"`
	Error  KrakenApiError `json:"error"`
}

/*
 Get tradable asset pairs.
 Note: If an asset pair is on a maker/taker fee schedule, the taker side
 is given in "fees" and maker side in "fees_maker".
 For pairs not on maker/taker, they will only be given in "fees".

 https://www.kraken.com/help/api#get-tradable-pairs
*/
func (k *Kraken) GetTradablePairs() (*AssetPairMap, error) {

	req, err := http.NewRequest("GET", URLGetTradablePairs, nil)
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
		return nil, errors.New("We got error" + dat.Error[0])
	}

	return &dat.Result, nil
}

type TickerInfo struct {
	/* Ask array(<price>, <whole lot volume>, <lot volume>). */
	A []string
	/* Bid array(<price>, <whole lot volume>, <lot volume>). */
	B []string
	/* Last trade closed array(<price>, <lot volume>). */
	C []string
	/* Volume array(<today>, <last 24 hours>). */
	V []string
	/* Volume weighted average price array(<today>, <last 24 hours>). */
	P []string
	/* Number of trades array(<today>, <last 24 hours>). */
	T []int
	/* Low array(<today>, <last 24 hours>). */
	L []string
	/* High array(<today>, <last 24 hours>). */
	H []string
	/* Today's opening price. */
	O string
}

type TickerInfoMap map[string]TickerInfo

type TickerInfoResult struct {
	Result TickerInfoMap  `json:"result"`
	Error  KrakenApiError `json:"error"`
}

/*
 Get ticker information.
 Input: comma delimited list of asset pairs to get info on
 Result: array of pair names and their ticker info

 https://www.kraken.com/help/api#get-ticker-info
*/
func (k *Kraken) GetTickerInfo(pairs []string) (*TickerInfoMap, error) {

	if len(pairs) == 0 {
		return nil, errors.New("Parameter pairs cannot be empty.")
	}
	var pairList string
	pairList = pairs[0]
	for i := 1; i < len(pairs); i++ {
		pairList += "," + pairs[i]
	}

	req, err := http.NewRequest("GET", URLGetTickerInfo, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("pair", pairList)
	req.URL.RawQuery = query.Encode()
	fmt.Println("We have URL that we are doing:    " + req.URL.String())
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
		return nil, errors.New("We got error: " + dat.Error[0])
	}

	return &dat.Result, nil
}
