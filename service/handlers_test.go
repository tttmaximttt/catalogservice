package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/codegangsta/negroni"
)

func TestGetDetailsForCatalogItemReturnsProperData(test *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := MakeTestServer()

	targetSKU := "THINGAMAJIG12"

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/catalog/"+targetSKU, nil)
	server.ServeHTTP(recorder, request)

	var detail catalogItem

	if recorder.Code != http.StatusOK {
		test.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
	}

	payload, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		test.Errorf("Error parsing response body: %v", err)
	}

	err = json.Unmarshal(payload, &detail)
	if err != nil {
		test.Errorf("Error unmarshaling response to catalog item: %v", err)
	}

	if detail.QuantityInStock != 1000 {
		test.Errorf("Expected 100 qty in stock, got %d", detail.QuantityInStock)
	}
	if detail.ShipsWithin != 99 {
		test.Errorf("Expected shipswithin 14 days, got %d", detail.ShipsWithin)
	}
	if detail.SKU != "THINGAMAJIG12" {
		test.Errorf("Expected SKU THINGAMAJIG12, got %s", detail.SKU)
	}
	if detail.ProductID != 1 {
		test.Errorf("Expected product ID of 1, got %d", detail.ProductID)
	}
}

func MakeTestServer() *negroni.Negroni {
	fakeClient := fakeWebClient{}
	return NewServerFromClient(fakeClient)
}

type fakeWebClient struct{}

func (client fakeWebClient) getFulfillmentStatus(sku string) (status fulfillmentStatus, err error) {
	status = fulfillmentStatus{
	SKU: sku,
		ShipsWithin: 99,
		QuantityInStock: 1000,
	}
	return status, err
}