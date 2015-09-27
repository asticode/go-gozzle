package gozzle

import (
	"bytes"
	"encoding/json"
	"github.com/asticode/go-parallelizator/parallelizator"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Gozzle struct {
	Configuration  Configuration
	parallelizator *parallelizator.Parallelizator
	client         Client
}

type Client struct {
	baseUrl     string
	maxSizeBody int
	httpClient  *http.Client
}

type Method int

const (
	METHOD_GET Method = iota
	METHOD_POST
	METHOD_PATCH
	METHOD_DELETE
)

var methodNames = map[Method]string{
	METHOD_GET:    "GET",
	METHOD_POST:   "POST",
	METHOD_PATCH:  "PATCH",
	METHOD_DELETE: "DELETE",
}

func NewGozzle(oConfiguration Configuration, oParallelizator *parallelizator.Parallelizator) (*Gozzle, error) {
	oGozzle := Gozzle{
		Configuration: oConfiguration,
	}
	oGozzle.new(oParallelizator)
	oErr := oGozzle.LoadConfiguration()
	return &oGozzle, oErr
}

func (oGozzle *Gozzle) LoadConfiguration() error {
	// Set base url
	oGozzle.client.baseUrl = oGozzle.Configuration.BaseUrl

	// Return
	return nil
}

func (oGozzle *Gozzle) new(oParallelizator *parallelizator.Parallelizator) {
	// Initialize client
	oGozzle.client.httpClient = &http.Client{}

	// Set parallelizator
	oGozzle.parallelizator = oParallelizator
}

func (oGozzle Gozzle) Exec(oRequestSet RequestSet) (*ResponseSet, error) {
	// Initialize
	oResponseSet := ResponseSet{}

	// Create wait group
	oWaitGroup := sync.WaitGroup{}
	oWaitGroup.Add(len(oRequestSet))

	// Loop through requests
	aErrors := []error{}
	for _, oRequest := range oRequestSet {
		// Exec request
		oGozzle.execRequestWrapper(oRequest, &oResponseSet, &aErrors, &oWaitGroup)
	}

	// Wait
	oWaitGroup.Wait()

	// Process errors
	for _, oErr := range aErrors {
		if oErr != nil {
			return &oResponseSet, oErr
		}
	}

	// Return
	return &oResponseSet, nil
}

func (oGozzle Gozzle) execRequestWrapper(oRequest *Request, oResponseSet *ResponseSet, aErrors *[]error, oWaitGroup *sync.WaitGroup) {
	oGozzle.parallelizator.AddJob(func() {
		// Exec
		*aErrors = append(*aErrors, oGozzle.execRequest(oRequest, oResponseSet))

		// Tell wait group it's done
		oWaitGroup.Done()
	})
}

func (oGozzle Gozzle) execRequest(oRequest *Request, oResponseSet *ResponseSet) error {
	// Before handler
	if oRequest.beforeHandler != nil {
		bContinue := oRequest.beforeHandler(oRequest)
		if bContinue != true {
			return nil
		}
	}

	// Get full url
	sFullUrl := oRequest.endpoint
	if oGozzle.client.baseUrl != "" {
		sFullUrl = oGozzle.client.baseUrl + sFullUrl
	}
	sFullUrl += "?"

	// Add query parameters
	for sQueryParameterKey, sQueryParameterValue := range oRequest.query {
		sFullUrl = sFullUrl + url.QueryEscape(sQueryParameterKey) + "=" + url.QueryEscape(sQueryParameterValue) + "&"
	}
	sFullUrl = strings.Trim(sFullUrl, "&")

	// Encode body
	oBody, oErr := json.Marshal(oRequest.body)
	if oErr != nil {
		return oErr
	}

	// Create http request
	oHttpRequest, oErr := http.NewRequest(methodNames[oRequest.method], sFullUrl, bytes.NewBuffer(oBody))
	if oErr != nil {
		return oErr
	}

	// Add headers.
	for sHeaderKey, sHeaderValue := range oRequest.headers {
		oHttpRequest.Header.Set(sHeaderKey, sHeaderValue)
	}

	// Send request
	oHttpResponse, oErr := oGozzle.client.httpClient.Do(oHttpRequest)
	if oErr != nil {
		return oErr
	}

	// Add response
	(*oResponseSet).addResponse(oRequest, oHttpResponse, oGozzle)

	// After handler
	if oRequest.afterHandler != nil {
		oRequest.afterHandler(oRequest, oResponseSet)
	}

	// Return
	return nil
}
