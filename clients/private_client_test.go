package clients

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"reflect"
	"testing"
	"time"
)

//
//
//

func Test_GetAccountReportStatus(t *testing.T) {
	// Setup the mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Mock the time request
	httpmock.RegisterResponder(
		"GET",
		"https://mock-api.gdax.com/reports/0428b97b-bec1-429e-a94c-59232926778d",
		httpmock.NewStringResponder(
			200,
			`
				{
					"id": "0428b97b-bec1-429e-a94c-59232926778d",
					"type": "fills",
					"status": "creating",
					"created_at": "2015-01-06T10:34:47.000Z",
					"expires_at": "2015-01-13T10:35:47.000Z",
					"params": {
						"start_date": "2014-11-01T00:00:00.000Z",
						"end_date": "2014-11-30T23:59:59.000Z"
					}
				}
			`,
		),
	)
	client := NewMockClient()
	output, err := GetAccountReportStatus(client, "0428b97b-bec1-429e-a94c-59232926778d")
	expected := AccountReportStatus{
		ID:          "0428b97b-bec1-429e-a94c-59232926778d",
		Type:        "fills",
		Status:      "creating",
		CreatedAt:   time.Date(2015, 01, 06, 10, 34, 47, 0, time.UTC),
		CompletedAt: time.Date(0001, 01, 01, 00, 00, 00, 0, time.UTC),
		ExpiresAt:   time.Date(2015, 01, 13, 10, 35, 47, 0, time.UTC),
		FileURL:     "",
	}
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Expected output to not be nil, actual = %v", output)
	}
	if !reflect.DeepEqual(output.ID, expected.ID) {
		t.Fatalf("Expected output.ID %v to match expected.ID %v", output.ID, expected.ID)
	}
	if !reflect.DeepEqual(output.Type, expected.Type) {
		t.Fatalf("Expected output.Type %v to match expected.Type %v", output.Type, expected.Type)
	}
	if !reflect.DeepEqual(output.Status, expected.Status) {
		t.Fatalf("Expected output.Status %v to match expected.Status %v", output.Status, expected.Status)
	}
	if !reflect.DeepEqual(output.CreatedAt.UTC().Unix(), expected.CreatedAt.UTC().Unix()) {
		t.Fatalf("Expected output.CreatedAt %v to match expected.CreatedAt %v", output.CreatedAt, expected.CreatedAt)
	}
	if !reflect.DeepEqual(output.CompletedAt.UTC().Unix(), expected.CompletedAt.UTC().Unix()) {
		t.Fatalf("Expected output.CompletedAt %v to match expected.CompletedAt %v", output.CompletedAt, expected.CompletedAt)
	}
	if !reflect.DeepEqual(output.FileURL, expected.FileURL) {
		t.Fatalf("Expected output.FileURL %v to match expected.FileURL %v", output.FileURL, expected.FileURL)
	}
}

//
//
//

func Test_mock_GetAccountTrailingVolume(t *testing.T) {
	// Setup the mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Mock the time request
	httpmock.RegisterResponder(
		"GET",
		"https://mock-api.gdax.com/users/self/trailing-volume",
		httpmock.NewStringResponder(
			200,
			`
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
			`,
		),
	)
	client := NewMockClient()
	expected := AccountTrailingVolume{
		ProductID:      "BTC-USD",
		ExchangeVolume: 11800.00000000,
		Volume:         100.00000000,
		RecordedAt:     time.Date(1973, 11, 29, 00, 05, 01, 123456*1000, time.UTC),
	}
	output, err := GetAccountTrailingVolume(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if len(output) != 2 {
		t.Fatalf("Expected output.length = 2, actual = %v", len(output))
	}
	if !reflect.DeepEqual(output[0].RecordedAt, expected.RecordedAt) {
		t.Fatalf("Expected output.RecordedAt %v to match expected.RecordedAt %v", output[0].RecordedAt, expected.RecordedAt)
	}
	if output[0].ProductID != expected.ProductID {
		t.Fatalf("Expected output.ProductID %v to match expected %v", output[0].ProductID, expected.ProductID)
	}
	if output[0].ExchangeVolume != expected.ExchangeVolume {
		t.Fatalf("Expected output.ExchangeVolume %v to match expected %v", output[0].ExchangeVolume, expected.ExchangeVolume)
	}
	if output[0].Volume != expected.Volume {
		t.Fatalf("Expected output.Volume %v to match expected %v", output[0].Volume, expected.Volume)
	}
}
