package clients

import (
	"net/url"
	"testing"
	"time"
)

//
//
//

func Test_live_GetProducts(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProducts(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) == 0 {
		t.Fatalf("Output should have more than item, %v", output)
	}
}

//
//
//

func Test_live_GetTime(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output := &GdaxTimeResponse{}
	err := client.Get("/time", url.Values{}, output)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
}

//
//
//

func Test_live_GetCurrencies(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetCurrencies(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) == 0 {
		t.Fatalf("Output should have >= 1 items, %v", output)
	}
}

//
//
//

func Test_live_GetProduct24HrStats(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProduct24HrStats(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
}

//
//
//

func Test_live_GetProductHistoricRates(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	start := time.Now().UTC().Add(-1 * HistoricRateGranularity_1day * time.Second)
	end := time.Now().UTC()
	output, err := GetProductHistoricRates(client, "BTC-USD", &start, &end, HistoricRateGranularity_1day)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
}

//
//
//

func Test_live_GetProductTrades(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProductTrades(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) == 0 {
		t.Fatalf("Output should have >= 1 items, %v", output)
	}
}

//
//
//

func Test_live_GetProductTicker(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProductTicker(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
}

//
//
//

func Test_live_GetProductOrderBookLevel1(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProductOrderBookLevel1(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
}

func Test_live_GetProductOrderBookLevel2(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProductOrderBookLevel2(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}

}

func Test_live_GetProductOrderBookLevel3(t *testing.T) {
	client := &HttpGdaxPublicClient{
		Domain: `api-public.sandbox.gdax.com`,
	}
	output, err := GetProductOrderBookLevel3(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
}

//
//
//
