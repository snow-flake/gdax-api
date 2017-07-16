package clients

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
	"reflect"
)

type mockGdaxPublicClient struct {
	json string
	err  error
}

func (client *mockGdaxPublicClient) Get(path string, params url.Values, target interface{}) error {
	if client.err != nil {
		return client.err
	}
	body := strings.NewReader(client.json)
	return json.NewDecoder(body).Decode(target)
}

func Test_mockGdaxPublicClient_Get_Success(t *testing.T) {
	type TestResponse struct {
		Iso   time.Time `json:"iso"`
		Epoch float64   `json:"epoch"`
	}
	client := &mockGdaxPublicClient{
		json: "{ \"iso\": \"2015-01-07T23:47:25.201Z\", \"epoch\": 1420674445.201 }",
		err:  nil,
	}
	output := &TestResponse{}
	err := client.Get("/time", url.Values{}, output)
	if nil != err {
		t.Fatal("Expected error to be nil, %v", err)
	}
}

func Test_mockGdaxPublicClient_Get_Error(t *testing.T) {
	type TestResponse struct {
		Iso   time.Time `json:"iso"`
		Epoch float64   `json:"epoch"`
	}
	client := &mockGdaxPublicClient{
		json: "",
		err:  fmt.Errorf("Example error"),
	}
	output := &TestResponse{}
	err := client.Get("/time", url.Values{}, output)
	if client.err != err {
		t.Fatal("Expected error to match inputs, %v", err)
	}
}

//
//
//

func Test_mock_GetProducts(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			[
				{
					"id": "BTC-USD",
					"base_currency": "BTC",
					"quote_currency": "USD",
					"base_min_size": "0.01",
					"base_max_size": "10000.00",
					"quote_increment": "0.01"
				}
			]
		`,
		err: nil,
	}
	output, err := GetProducts(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) != 1 {
		t.Fatalf("Output should have 1 item, %v", output)
	}
	item := (output)[0]
	if item.ID != "BTC-USD" {
		t.Fatalf("Expected item.id = BTC-USD, actual = %v", item.ID)
	}
	if item.BaseCurrency != "BTC" {
		t.Fatalf("Expected item.base_currency = BTC, actual = %v", item.BaseCurrency)
	}
	if item.QuoteCurrency != "USD" {
		t.Fatalf("Expected item.quote_currency = USD, actual = %v", item.QuoteCurrency)
	}
	if item.BaseMinSize != "0.01" {
		t.Fatalf("Expected item.base_min_size = 0.01, actual = %v", item.BaseMinSize)
	}
	if item.BaseMaxSize != "10000.00" {
		t.Fatalf("Expected item.base_max_size = 10000.00, actual = %v", item.BaseMaxSize)
	}
	if item.QuoteIncrement != "0.01" {
		t.Fatalf("Expected item.quote_increment = 0.01, actual = %v", item.QuoteIncrement)
	}
}

//
//
//

func Test_mock_GetTime(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
				"iso": "2015-01-07T23:47:25.201Z",
				"epoch": 1420674445.201
			}
		`,
		err: nil,
	}
	expected := &GdaxTimeResponse{
		Iso: time.Date(2015, 01, 07, 23, 47, 25, 000000201, time.UTC),
		Epoch: 1420674445.201,
	}
	output, err := GetTime(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if !reflect.DeepEqual(output.Epoch, expected.Epoch) {
		t.Fatalf("Expected output.Epoch %v to match expected.Epoch %v", output.Epoch, expected.Epoch)
	}
	if !reflect.DeepEqual(output.Iso.Unix(), expected.Iso.Unix()) {
		t.Fatalf("Expected output %v to match expected %v", output, expected)
	}
}

//
//
//

func Test_mock_GetCurrencies(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			[
				{ "id": "BTC", "name": "Bitcoin", "min_size": "0.00000001" },
				{ "id": "USD", "name": "United States Dollar", "min_size": "0.01000000" }
			]
		`,
		err: nil,
	}
	expected := []GdaxCurrency{
		GdaxCurrency{ ID: "BTC", Name: "Bitcoin", MinSize: 0.00000001},
		GdaxCurrency{ ID: "USD", Name: "United States Dollar", MinSize: 0.01000000},

	}
	output, err := GetCurrencies(client)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) != 2 {
		t.Fatalf("Output should have 1 item, %v", output)
	}
	for i, actual := range(output) {
	if !reflect.DeepEqual(actual, expected[i]) {
		t.Fatalf("Expected output %v to match expected %v", actual, expected[i])
	}}
}

//
//
//

func Test_mock_GetProduct24HrStats(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
				"open":"2000.00000000",
				"high":"2110.06000000",
				"low":"1758.20000000",
				"volume":"20465.01966891",
				"last":"1893.91000000",
				"volume_30day":"398368.6657624"
			}
		`,
		err: nil,
	}
	expected := &GdaxProduct24HrStatsResponse{
		Open: 2000.00000000,
		High: 2110.06000000,
		Low: 1758.20000000,
		Volume: 20465.01966891,
		Last: 1893.91000000,
		Volume30Day: 398368.6657624,

	}
	output, err := GetProduct24HrStats(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if !reflect.DeepEqual(output, expected) {
		t.Fatalf("Expected output %v to match expected %v", output, expected)
	}
}

//
//
//

func Test_mock_GetProductHistoricRates(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			[
			  [1500130020,181.8,181.81,181.8,181.81,11.34496359],
			  [1500130005,181.81,181.81,181.81,181.81,5.75798592]
			]
		`,
		err: nil,
	}
	start := time.Now().UTC().Add(-1 * HistoricRateGranularity_5m * time.Second)
	end := time.Now().UTC()
	output, err := GetProductHistoricRates(client, "BTC-USD", &start, &end, HistoricRateGranularity_5m)
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) != 2 {
		t.Fatalf("Output should have 1 item, %v", output)
	}
	item := (output)[0]
	if item.Time != time.Unix(1500130020, 0) {
		t.Fatalf("Expected item.Time = 1500130020, actual = %v", item.Time)
	}
	if item.Low != 181.8 {
		t.Fatalf("Expected item.Low = 181.8, actual = %v", item.Low)
	}
	if item.High != 181.81 {
		t.Fatalf("Expected item.High = 181.81, actual = %v", item.High)
	}
	if item.Open != 181.8 {
		t.Fatalf("Expected item.Open = 181.8, actual = %v", item.Open)
	}
	if item.Close != 181.81 {
		t.Fatalf("Expected item.Close = 181.81, actual = %v", item.Close)
	}
	if item.Volume != 11.34496359 {
		t.Fatalf("Expected item.Volume = 11.34496359, actual = %v", item.Volume)
	}
}

//
//
//

func Test_mock_GetProductTrades(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			[{
			    "time": "2014-11-07T22:19:28.578544Z",
			    "trade_id": 74,
			    "price": "10.00000000",
			    "size": "0.01000000",
			    "side": "buy"
			}, {
			    "time": "2014-11-07T01:08:43.642366Z",
			    "trade_id": 73,
			    "price": "100.00000000",
			    "size": "0.01000000",
			    "side": "sell"
			}]
		`,
		err: nil,
	}
	output, err := GetProductTrades(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if len(output) != 2 {
		t.Fatalf("Output should have 1 item, %v", output)
	}
	item := (output)[0]
	expectedTime := time.Date(2014, 11, 07, 22, 19, 28, 578544, time.UTC).Unix()
	if item.Time.Unix() != expectedTime {
		t.Fatalf("Expected item.Time = 2014-11-07T22:19:28.578544Z, actual = %v", item.Time)
	}
	if item.TradeID != 74 {
		t.Fatalf("Expected item.TradeID = 74, actual = %v", item.TradeID)
	}
	if item.Price != 10.00000000 {
		t.Fatalf("Expected item.Price = 10.00000000, actual = %v", item.Price)
	}
	if item.Size != 0.01000000 {
		t.Fatalf("Expected item.Size = 0.01000000, actual = %v", item.Size)
	}
	if item.Side != "buy" {
		t.Fatalf("Expected item.Side = buy, actual = %v", item.Side)
	}
}

//
//
//

func Test_mock_GetProductTicker(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
			  "trade_id": 4729088,
			  "price": "333.99",
			  "size": "0.193",
			  "bid": "333.98",
			  "ask": "333.99",
			  "volume": "5957.11914015",
			  "time": "2015-11-14T20:46:03.511254Z"
			}
		`,
		err: nil,
	}
	expected := &GdaxProductTickerResponse{
		TradeID: 4729088,
		Price: 333.99,
		Size: 0.193,
		Bid: 333.98,
		Ask: 333.99,
		Volume: 5957.11914015,
		Time: time.Date(2015,11,14,20,46,03,511254, time.UTC),
	}
	output, err := GetProductTicker(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	expectedTime := time.Date(2015, 11, 14, 20, 46, 03, 511254, time.UTC).Unix()
	if output.Time.Unix() != expectedTime {
		t.Fatalf("Expected item.Time = 2015-11-14T20:46:03.511254Z, actual = %v", output.Time)
	}
	if output.TradeID != expected.TradeID {
		t.Fatalf("Expected output.TradeId = 4729088, actual = %v", output.TradeID)
	}
	if output.Price != expected.Price {
		t.Fatalf("Expected output.Price = 333.99, actual = %v", output.Price)
	}
	if output.Size != expected.Size {
		t.Fatalf("Expected output.Size = 0.193, actual = %v", output.Size)
	}
	if output.Bid != expected.Bid {
		t.Fatalf("Expected output.Bid = 333.98, actual = %v", output.Bid)
	}
	if output.Ask != expected.Ask {
		t.Fatalf("Expected output.Ask = 333.99, actual = %v", output.Ask)
	}
	if output.Volume != expected.Volume {
		t.Fatalf("Expected output.Volume = 5957.11914015, actual = %v", output.Volume)
	}
}

//
//
//

func Test_mock_GetProductOrderBookLevel1(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
				"sequence":775966773,
				"bids":[["180.79","142.55091057",2]],
				"asks":[["180.84","9.91691592",2]]
			}
		`,
		err: nil,
	}
	output, err := GetProductOrderBookLevel1(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if output.Sequence != 775966773 {
		t.Fatalf("Expected output.Sequence = 775966773, actual = %v", output.Sequence)
	}
	if output.Bid.Price != 180.79 {
		t.Fatalf("Expected output.Bid.Price = 180.79, actual = %v", output.Bid.Price)
	}
	if output.Bid.Size != 142.55091057 {
		t.Fatalf("Expected output.Bid.Size = 142.55091057, actual = %v", output.Bid.Size)
	}
	if output.Bid.NumOrders != 2 {
		t.Fatalf("Expected output.Bid.NumOrders = 2, actual = %v", output.Bid.NumOrders)
	}
	if output.Ask.Price != 180.84 {
		t.Fatalf("Expected output.Ask.Price = 180.84, actual = %v", output.Ask.Price)
	}
	if output.Ask.Size != 9.91691592 {
		t.Fatalf("Expected output.Ask.Size = 9.91691592, actual = %v", output.Ask.Size)
	}
	if output.Ask.NumOrders != 2 {
		t.Fatalf("Expected output.Ask.NumOrders = 2, actual = %v", output.Ask.NumOrders)
	}
}

func Test_mock_GetProductOrderBookLevel2(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
				"sequence":776000158,
				"bids":[
					["179.32","45.346",4],
					["179.29","4",1],
					["179.21","4",1],
					["179.13","4.2",1],
					["179.04","3",1],
					["179.02","1",1],
					["179","1",1],
					["178.97","3",1],
					["178.96","4.61",2],
					["178.91","1",1],
					["178.89","4",2],
					["178.88","6",3],
					["178.8","3.4",1],
					["178.78","1",1],
					["178.73","4",1],
					["178.64","5",2],
					["178.62","1",1],
					["178.57","4",1],
					["178.53","4",1],
					["178.5","20.584",2],
					["178.49","26.67193",1],
					["178.48","3.9",1],
					["178.42","17.3144",1],
					["178.41","3.8",1],
					["178.4","1",1],
					["178.35","48.05",1],
					["178.32","52.52",2],
					["178.27","55.30032",1],
					["178.26","53.14",1],
					["178.25","4",1],
					["178.22","92.89",2],
					["178.21","45.666",2],
					["178.2","144.86150354",1],
					["178.17","4",1],
					["178.15","1",1],
					["178.11","1.1",2],
					["178.1","11.2135",4],
					["178.09","0.63024002",1],
					["178.08","4",1],
					["178.06","50",1],
					["178.02","3.759",3],
					["178.01","0.1",1],
					["178","206.72277694",9],
					["177.97","75",1],
					["177.96","1",1],
					["177.95","2",2],
					["177.92","3.22670157",2],
					["177.91","4.1",2],
					["177.85","10",1],
					["177.83","3.4",1]
				],
				"asks":[["179.33","201.26000475",8],
					["179.34","313.63117673",2],
					["179.46","0.25",1],
					["179.5","2.7",2],
					["179.51","17.81115125",3],
					["179.52","11.799",2],
					["179.53","0.28",3],
					["179.56","26.07368941",1],
					["179.75","6",2],
					["179.76","4",1],
					["179.8","0.01",1],
					["179.83","4.3",1],
					["179.89","0.01",1],
					["179.91","3.8",1],
					["179.93","1.01",2],
					["179.94","1",1],
					["179.99","71.33645493",2],
					["180","739.06974762",4],
					["180.07","4",1],
					["180.08","0.01",1],
					["180.16","4.46",2],
					["180.17","1",1],
					["180.19","1",1],
					["180.23","4",1],
					["180.24","1",1],
					["180.28","2.3",2],
					["180.29","0.01",1],
					["180.31","4.3",1],
					["180.38","1",1],
					["180.4","3",1],
					["180.43","1.876",2],
					["180.44","25.1981",1],
					["180.46","1.01",2],
					["180.47","4",1],
					["180.49","0.1",1],
					["180.54","1",1],
					["180.55","4",1],
					["180.58","3",1],
					["180.59","1.1",2],
					["180.62","4",1],
					["180.63","3",1],
					["180.66","18.3144",2],
					["180.67","1",1],
					["180.68","1.068",2],
					["180.69","0.1",1],
					["180.71","3.4",1],
					["180.79","10.4331883",3],
					["180.8","12.5",4],
					["180.88","3",1],
					["180.89","0.1",1]
				]
			}
		`,
		err: nil,
	}
	output, err := GetProductOrderBookLevel2(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if output.Sequence != 776000158 {
		t.Fatalf("Expected output.Sequence = 776000158, actual = %v", output.Sequence)
	}
	if len(output.Bids) != 50 {
		t.Fatalf("Expected output.Bids.length = 50, actual = %v", len(output.Bids))
	}
	if len(output.Asks) != 50 {
		t.Fatalf("Expected output.Asks.length = 50, actual = %v", len(output.Asks))
	}
	if output.Bids[0].Price != 179.32 {
		t.Fatalf("Expected output.Bid.Price = 179.32, actual = %v", output.Bids[0].Price)
	}
	if output.Bids[0].Size != 45.346 {
		t.Fatalf("Expected output.Bid.Size = 45.346, actual = %v", output.Bids[0].Size)
	}
	if output.Bids[0].NumOrders != 4 {
		t.Fatalf("Expected output.Bid.NumOrders = 4, actual = %v", output.Bids[0].NumOrders)
	}
	if output.Asks[0].Price != 179.33 {
		t.Fatalf("Expected output.Ask.Price = 179.33, actual = %v", output.Asks[0].Price)
	}
	if output.Asks[0].Size != 201.26000475 {
		t.Fatalf("Expected output.Ask.Size = 201.26000475, actual = %v", output.Asks[0].Size)
	}
	if output.Asks[0].NumOrders != 8 {
		t.Fatalf("Expected output.Ask.NumOrders = 8, actual = %v", output.Asks[0].NumOrders)
	}
}

func Test_mock_GetProductOrderBookLevel3(t *testing.T) {
	client := &mockGdaxPublicClient{
		json: `
			{
				"sequence":776035515,
				"bids":[
					["179.27","0.27","b2276903-d242-445d-bc86-6ca2f65f5ed7"],
					["179.27","24.91831786","5235f930-413a-4e69-a6c7-4a0068e75da7"],
					["179.27","54","247a2eee-d2c4-4b20-a37a-3ad150d8193f"]
				],
				"asks":[
					["179.28","252.43268514","e2a982f2-5cd0-4775-ab36-08f79d622e2b"],
					["179.36","4.5","afda1fa9-e355-4056-a3b8-1686806dcc73"],
					["179.38","0.25","b9f8e2e2-4bc5-4ff5-a024-f6b737ad19e7"]
				]
			}
		`,
		err: nil,
	}
	output, err := GetProductOrderBookLevel3(client, "BTC-USD")
	if err != nil {
		t.Fatalf("Error should be nil, %v", err)
	}
	if output == nil {
		t.Fatalf("Output should not be nil, %v", output)
	}
	if output.Sequence != 776035515 {
		t.Fatalf("Expected output.Sequence = 776035515, actual = %v", output.Sequence)
	}
	if len(output.Bids) != 3 {
		t.Fatalf("Expected output.Bids.length = 50, actual = %v", len(output.Bids))
	}
	if len(output.Asks) != 3 {
		t.Fatalf("Expected output.Asks.length = 50, actual = %v", len(output.Asks))
	}
	if output.Bids[0].Price != 179.27 {
		t.Fatalf("Expected output.Bid.Price = 179.27, actual = %v", output.Bids[0].Price)
	}
	if output.Bids[0].Size != 0.27 {
		t.Fatalf("Expected output.Bid.Size = 0.27, actual = %v", output.Bids[0].Size)
	}
	if output.Bids[0].OrderId != "b2276903-d242-445d-bc86-6ca2f65f5ed7" {
		t.Fatalf("Expected output.Bid.OrderId = b2276903-d242-445d-bc86-6ca2f65f5ed7, actual = %v", output.Bids[0].OrderId)
	}
	if output.Asks[0].Price != 179.28 {
		t.Fatalf("Expected output.Ask.Price = 179.28, actual = %v", output.Asks[0].Price)
	}
	if output.Asks[0].Size != 252.43268514 {
		t.Fatalf("Expected output.Ask.Size = 252.43268514, actual = %v", output.Asks[0].Size)
	}
	if output.Asks[0].OrderId != "e2a982f2-5cd0-4775-ab36-08f79d622e2b" {
		t.Fatalf("Expected output.Ask.OrderId = e2a982f2-5cd0-4775-ab36-08f79d622e2b, actual = %v", output.Asks[0].OrderId)
	}
}

//
//
//
