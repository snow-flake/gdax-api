package clients

import (
	"net/url"
	"time"
)

type AccountTrailingVolume struct {
	ProductID      string    `json:"product_id"`
	ExchangeVolume float64   `json:"exchange_volume,string"`
	Volume         float64   `json:"volume,string"`
	RecordedAt     time.Time `json:"recorded_at"`
}
type AccountTrailingVolumeResponse []AccountTrailingVolume

/*
User Account: Trailing Volume

HTTP REQUEST
GET /users/self/trailing-volume

This request will return your 30-day trailing volume for all products. This is a cached value that’s calculated every day at midnight UTC.

HTTP RESPONSE
Trailing Volume
[
	{
		"product_id": "BTC-USD",
		"exchange_volume": "11800.00000000",
		"volume": "100.00000000",
		"recorded_at": "1973-11-29T00:05:01.123456Z"
	},
	{
		"product_id": "LTC-USD",
		"exchange_volume": "51010.04100000",
		"volume": "2010.04100000",
		"recorded_at": "1973-11-29T00:05:02.123456Z"
	}
]
*/
func GetAccountTrailingVolume(client *Client) (AccountTrailingVolumeResponse, error) {
	pathname := "/users/self/trailing-volume"
	params := url.Values{}
	output := []AccountTrailingVolume{}
	_, err := client.Get(pathname, params, &output)
	return output, err
}
