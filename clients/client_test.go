package clients

import (
	"net/url"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"reflect"
)

func Test_Client_Get(t *testing.T) {
	// Setup the mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Mock the time request
	httpmock.RegisterResponder(
		"GET",
		"https://mock-api.gdax.com/time",
		httpmock.NewStringResponder(
			200,
			`{ "iso": "2015-01-07T23:47:25.201Z", "epoch": 1420674445.201 }`,
		),
	)
	// Create the variables for the test API request
	client := NewMockClient()
	pathname := "/time"
	url_params := url.Values{}
	output := &GdaxTimeResponse{}
	expected := &GdaxTimeResponse{
		Iso:   time.Date(2015, 01, 07, 23, 47, 25, 000000201, time.UTC),
		Epoch: 1420674445.201,
	}
	// Execute the stub'd request
	res, err := client.Get(pathname, url_params, output)
	if nil != err {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if nil == res {
		t.Fatalf("Expected response to not be nil, actual = %v", res)
	}
	if !reflect.DeepEqual(output.Epoch, expected.Epoch) {
		t.Fatalf("Expected output.Epoch %v to match expected.Epoch %v", output.Epoch, expected.Epoch)
	}
	if !reflect.DeepEqual(output.Iso.Unix(), expected.Iso.Unix()) {
		t.Fatalf("Expected output %v to match expected %v", output, expected)
	}
}

func Test_Client_Post(t *testing.T) {
	// Setup the mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Mock the time request
	httpmock.RegisterResponder(
		"POST",
		"https://mock-api.gdax.com/time",
		httpmock.NewStringResponder(
			200,
			`{ "iso": "2015-01-07T23:47:25.201Z", "epoch": 1420674445.201 }`,
		),
	)
	// Create the variables for the test API request
	client := NewMockClient()
	pathname := "/time"
	url_params := map[string]string{
		"test": "value",
	}
	output := &GdaxTimeResponse{}
	expected := &GdaxTimeResponse{
		Iso:   time.Date(2015, 01, 07, 23, 47, 25, 000000201, time.UTC),
		Epoch: 1420674445.201,
	}
	// Execute the stub'd request
	res, err := client.Post(pathname, url_params, output)
	if nil != err {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if nil == res {
		t.Fatalf("Expected response to not be nil, actual = %v", res)
	}
	if !reflect.DeepEqual(output.Epoch, expected.Epoch) {
		t.Fatalf("Expected output.Epoch %v to match expected.Epoch %v", output.Epoch, expected.Epoch)
	}
	if !reflect.DeepEqual(output.Iso.Unix(), expected.Iso.Unix()) {
		t.Fatalf("Expected output %v to match expected %v", output, expected)
	}
}

func Test_Client_Delete(t *testing.T) {
	// Setup the mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Mock the time request
	httpmock.RegisterResponder(
		"DELETE",
		"https://mock-api.gdax.com/time",
		httpmock.NewStringResponder(
			200,
			`{ "iso": "2015-01-07T23:47:25.201Z", "epoch": 1420674445.201 }`,
		),
	)
	// Create the variables for the test API request
	client := NewMockClient()
	pathname := "/time"
	url_params := map[string]string{
		"test": "value",
	}
	output := &GdaxTimeResponse{}
	expected := &GdaxTimeResponse{
		Iso:   time.Date(2015, 01, 07, 23, 47, 25, 000000201, time.UTC),
		Epoch: 1420674445.201,
	}
	// Execute the stub'd request
	res, err := client.Delete(pathname, url_params, output)
	if nil != err {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if nil == res {
		t.Fatalf("Expected response to not be nil, actual = %v", res)
	}
	if !reflect.DeepEqual(output.Epoch, expected.Epoch) {
		t.Fatalf("Expected output.Epoch %v to match expected.Epoch %v", output.Epoch, expected.Epoch)
	}
	if !reflect.DeepEqual(output.Iso.Unix(), expected.Iso.Unix()) {
		t.Fatalf("Expected output %v to match expected %v", output, expected)
	}
}

func Test_Client_formatUrl(t *testing.T) {
	client := NewMockClient()
	output := client.formatUrl(
		"/test",
		url.Values{
			"test": []string{"value"},
		},
	)
	if output != "/test?test=value" {
		t.Fatalf("Expected url to match /test?test=value, actual = %v", output)
	}
}

func Test_Client_encodeBody(t *testing.T) {
	client := NewMockClient()
	reader, data, err := client.encodeBody(map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if reader == nil {
		t.Fatalf("Expected reader to not be nil, actual = %v", reader)
	}
	if string(data) != `{"key":"value"}` {
		t.Fatalf("Expected data to match data, actual = %v", string(data))
	}
}

func Test_Client_generateMessageSignature(t *testing.T) {
	client := NewMockClient()
	timestamp := "2014-11-06T10:34:47.123456Z"
	method := "GET"
	partial_url := "/time?test=params"
	encoded_data := []byte(`{"test": "body"}`)
	signature, err := client.generateMessageSignature(timestamp, method, partial_url, encoded_data)
	if nil != err {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if signature == "" {
		t.Fatalf(signature)
	}
}

func Test_Client_generateMessageSignature_isBlank(t *testing.T) {
	client := NewSandboxClient("", "", "")
	timestamp := "2014-11-06T10:34:47.123456Z"
	method := "GET"
	partial_url := "/time?test=params"
	encoded_data := []byte(`{"test": "body"}`)
	signature, err := client.generateMessageSignature(timestamp, method, partial_url, encoded_data)
	if nil != err {
		t.Fatalf("Expected error to be nil, actual = %v", err)
	}
	if signature != "" {
		t.Fatalf(signature)
	}
}

func Test_Client_decodeError(t *testing.T) {
	body_data := []byte(`
		{
			"message": "Invalid Price"
		}
	`)
	client := NewSandboxClient("", "", "")
	err := client.decodeError(body_data)
	if nil == err {
		t.Fatalf("Expected to return error, actual = %v", err)
	}
	if err.Error() != "Invalid Price" {
		t.Fatalf("Expected to return error, actual = %v", err)
	}
}
