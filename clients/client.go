package clients

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
	API Client for Coinbase GDAX

	Sandbox Endpoints:
		Website: https://public.sandbox.gdax.com
		REST API: https://api-public.sandbox.gdax.com
		Websocket Feed: wss://ws-feed-public.sandbox.gdax.com
		FIX API: https://fix-public.sandbox.gdax.com

	Production Endpoints:
		Website: https://gdax.com
		REST API: https://api.gdax.com
		Websocket Feed: wss://ws-feed.gdax.com
		FIX API: tcp+ssl://fix.gdax.com:4198
*/
type Client struct {
	URL        string
	Secret     string
	Key        string
	Passphrase string
}

type ClientError struct {
	Message string `json:"message"`
}

func (e ClientError) Error() string {
	return e.Message
}

func NewProductionClient(secret, key, passphrase string) *Client {
	return &Client{
		URL:        "https://api.gdax.com",
		Secret:     secret,
		Key:        key,
		Passphrase: passphrase,
	}
}

func NewSandboxClient(secret, key, passphrase string) *Client {
	return &Client{
		URL:        "https://api-public.sandbox.gdax.com",
		Secret:     secret,
		Key:        key,
		Passphrase: passphrase,
	}
}

func NewMockClient() *Client {
	secret := "c3VwZXItc2VjcmV0LXBhc3N3b3Jk"         // super-secret-password
	key := "YW1hemluZy1zdXBlci1zZWNyZXQta2V5"        // amazing-super-secret-key
	passphrase := "YW1hemluZy1zdXBlci1wYXNzcGhyYXNl" // amazing-super-passphrase
	return &Client{
		URL:        "https://mock-api.gdax.com",
		Secret:     secret,
		Key:        key,
		Passphrase: passphrase,
	}
}

func (c *Client) Get(pathname string, url_params url.Values, result interface{}) (res *http.Response, err error) {
	return c.request("GET", pathname, url_params, nil, result)
}

func (c *Client) Post(pathname string, body_params, result interface{}) (res *http.Response, err error) {
	return c.request("POST", pathname, nil, body_params, result)
}

func (c *Client) Delete(pathname string, body_params, result interface{}) (res *http.Response, err error) {
	return c.request("DELETE", pathname, nil, body_params, result)
}

/*
	Requests:
		All requests and responses are application/json content type and
		follow typical HTTP response status codes for success and failure.

	Errors:
		Unless otherwise stated, errors to bad requests will respond with HTTP 4xx or status codes.
		The body will also contain a message parameter indicating the cause.
		Your language’s http library should be configured to provide message bodies for non-2xx requests
		so that you can read the message field from the body.
	{
		"message": "Invalid Price"
	}

	Common error codes:
	| Status Code | Reason                                                       |
	| 400         | Bad Request – Invalid request format                         |
	| 401         | Unauthorized – Invalid API Key                               |
	| 403         | Forbidden – You do not have access to the requested resource |
	| 404         | Not Found                                                    |
	| 500         | Internal Server Error – We had a problem with our server     |
*/
func (c *Client) request(method string, pathname string, url_params url.Values, body_params, result interface{}) (*http.Response, error) {
	// Generate the current timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// Format the url with "/pathname?query=params"
	partial_url := c.formatUrl(pathname, url_params)
	// Encode the message body as a JSON blob
	body, encoded_data, err := c.encodeBody(body_params)
	if err != nil {
		return nil, err
	}
	// Generate the message signature
	signature, err := c.generateMessageSignature(timestamp, method, partial_url, encoded_data)
	if err != nil {
		return nil, err
	}
	// Finally create the HTTP request with the given url, body, and headers
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.URL, partial_url), body)
	if err != nil {
		return nil, err
	}
	// Add the headers to the request
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla")
	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	if "" != c.Key {
		req.Header.Add("CB-ACCESS-KEY", c.Key)
	}
	if "" != c.Passphrase {
		req.Header.Add("CB-ACCESS-PASSPHRASE", c.Passphrase)
	}
	if "" != signature {
		req.Header.Add("CB-ACCESS-SIGN", signature)
	}
	// Execute the HTTP request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return res, err
	}
	// Read in the response body
	defer res.Body.Close()
	body_data, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return res, err
	}

	log.Printf("Request URL: %v", partial_url)
	log.Printf("Request Body: %v", string(body_data))
	log.Printf("Response Status: %v", res.StatusCode)
	log.Printf("Response Body: %v", string(body_data))

	// If the status code is !== 200 then bail now
	if res.StatusCode != 200 {
		return res, c.decodeError(body_data)
	}
	// Decode the body and return the output
	err = json.NewDecoder(bytes.NewReader(body_data)).Decode(result)
	return res, err
}

/*
	Format the full URL of the request
*/
func (c *Client) formatUrl(pathname string, params url.Values) string {
	encoded_params := url.Values(params).Encode()
	if encoded_params == "" {
		return pathname
	}
	return fmt.Sprintf("%s?%s", pathname, encoded_params)
}

/*
	Encode the body as a JSON encoded string
*/
func (c *Client) encodeBody(params interface{}) (*bytes.Reader, []byte, error) {
	body := bytes.NewReader(make([]byte, 0))
	if params == nil {
		return body, []byte{}, nil
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, data, err
	}

	body = bytes.NewReader(data)
	return body, data, nil

}

/*
 Generate the client signature for the request

 The CB-ACCESS-SIGN header is generated by creating a sha256 HMAC using the base64-decoded secret key on the prehash
 string timestamp + method + requestPath + body (where + represents string concatenation) and base64-encode the output.

 The timestamp value is the same as the CB-ACCESS-TIMESTAMP header.

 The body is the request body string or omitted if there is no request body (typically for GET requests).
 The method should be UPPER CASE.
*/
func (c *Client) generateMessageSignature(timestamp, method, partial_url string, encoded_data []byte) (string, error) {
	if c.Secret == "" {
		return "", nil
	}
	// Decode the secret key
	key, err := base64.StdEncoding.DecodeString(c.Secret)
	if err != nil {
		return "", err
	}
	// Format the message body
	message := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(method), partial_url, string(encoded_data))
	// Sign the message body
	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", err
	}
	// Encode the signed message into a base64 string
	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}

/*
 Decode the error response
*/
func (c *Client) decodeError(body_data []byte) error {
	client_error := ClientError{}
	reader := bytes.NewReader(body_data)
	err := json.NewDecoder(reader).Decode(&client_error)
	if err != nil {
		return err
	}
	return error(client_error)

}
