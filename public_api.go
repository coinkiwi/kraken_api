package kraken_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URLGetServerTime string = "https://api.kraken.com/0/public/Time"
const URLGetAssetsInfo string = "https://api.kraken.com/0/public/Assets"
const URLGetTradablePairs string = "https://api.kraken.com/0/public/AssetPairs"

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

	return &dat.Result, nil
}
