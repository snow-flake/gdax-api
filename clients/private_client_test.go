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
