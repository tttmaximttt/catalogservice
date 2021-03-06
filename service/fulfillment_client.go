package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type fulfillmentClient interface {
	getFulfillmentStatus(sku string) (status fulfillmentStatus, err error)
}

type fulfillmentWebClient struct {
	rootURL string
	skus string
}

func (client fulfillmentWebClient) getFulfillmentStatus(sku string) (status fulfillmentStatus, err error) {
	httpclient := &http.Client{}

	skuURL := fmt.Sprintf("%s/%s", client.rootURL + client.skus, sku)
	fmt.Printf("About to request SKU details from backing service: %s\n", skuURL)
	req, _ := http.NewRequest("GET", skuURL, nil)

	resp, err := httpclient.Do(req)

	if err != nil {
		fmt.Printf("Errored when sending request to the server: %s\n", err.Error())
		return
	}

	defer resp.Body.Close()
	payload, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(payload, &status)
	if err != nil {
		fmt.Println("Failed to unmarshal server response.")
		return
	}

	return status, err
}