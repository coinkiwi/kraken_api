package kraken_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URLGetServerTime string = "https://api.kraken.com/0/public/Time"
const URLGetAssetsInfo string = "https://api.kraken.com/0/public/Assets"

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
