package service

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/hudl/fargo"
)

func NewServerFromCFEnv(appEnv *cfenv.App) *negroni.Negroni {
	discovery := fargo.NewConn("http://localhost:8080/eureka/v2")
	app, _ := discovery.GetApp("FULLFILLMENT_APP")

	fmt.Println(app.Instances[0].IPAddr, app.Instances[0].Port)
	webClient := fulfillmentWebClient{
		rootURL: fmt.Sprintf("http://%s:%v/", app.Instances[0].IPAddr, app.Instances[0].Port),
		skus: "skus",
	}

	fmt.Println(webClient.rootURL)

	val, err := cftools.GetVCAPServiceProperty("backing-fulfill", "url", appEnv)
	if err == nil {
		webClient.rootURL = val
	} else {
		fmt.Printf("Failed to get URL property from bound service: %v\n", err)
	}
	fmt.Printf("Using the following URL for fulfillment backing service: %s\n", webClient.rootURL)

	return NewServerFromClient(webClient)
}

// NewServerFromClient configures and returns a Server.
func NewServerFromClient(webClient fulfillmentClient) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, webClient)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, webClient fulfillmentClient) {
	mx.HandleFunc("/", rootHandler(formatter)).Methods("GET")
	mx.HandleFunc("/catalog", getAllCatalogItemsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/catalog/{sku}", getCatalogItemDetailsHandler(formatter, webClient)).Methods("GET")
}