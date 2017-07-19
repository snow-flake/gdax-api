package clients

import (
	"fmt"
	"net/url"
	"time"
)

type AccountReportStatus struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	FileURL     string    `json:"file_url"`
}

/*
	Get report status
	HTTP REQUEST
		GET /reports/:report_id

	HTTP RESPONSE
	Response (creating report)
	{
		"id": "0428b97b-bec1-429e-a94c-59232926778d",
		"type": "fills",
		"status": "creating",
		"created_at": "2015-01-06T10:34:47.000Z",
		"completed_at": undefined,
		"expires_at": "2015-01-13T10:35:47.000Z",
		"file_url": undefined,
		"params": {
			"start_date": "2014-11-01T00:00:00.000Z",
			"end_date": "2014-11-30T23:59:59.000Z"
		}
	}
	Response (finished report)
	{
		"id": "0428b97b-bec1-429e-a94c-59232926778d",
		"type": "fills",
		"status": "ready",
		"created_at": "2015-01-06T10:34:47.000Z",
		"completed_at": "2015-01-06T10:35:47.000Z",
		"expires_at": "2015-01-13T10:35:47.000Z",
		"file_url": "https://example.com/0428b97b.../fills.pdf",
		"params": {
			"start_date": "2014-11-01T00:00:00.000Z",
			"end_date": "2014-11-30T23:59:59.000Z"
		}
	}

	Once a report request has been accepted for processing, the status is available by polling the report resource endpoint.
	The final report will be uploaded and available at file_url once the status indicates ready

	STATUS
	| Status | Description |
	| pending | The report request has been accepted and is awaiting processing |
	| creating | The report is being created |
	| ready | The report is ready for download from file_url |
*/
func GetAccountReportStatus(client *Client, report_id string) (*AccountReportStatus, error) {
	pathname := fmt.Sprintf("/reports/%s", report_id)
	params := url.Values{}
	output := &AccountReportStatus{}
	_, err := client.Get(pathname, params, &output)
	return output, err
}

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
