package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type IGdaxPrivateClient interface {
	Get(path string, params url.Values, target interface{}) error
}

type HttpGdaxPrivateClient struct {
	Domain string
}

func (client *HttpGdaxPrivateClient) Get(path string, params url.Values, target interface{}) error {
	url := fmt.Sprintf("https://%s%s?%s", client.Domain, path, url.Values(params).Encode())
	log.Println(url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
