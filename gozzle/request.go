package gozzle
import (
	"github.com/asticode/go-toolbox/array"
)

type Request struct {
	name          string
	method        Method
	endpoint      string
	headers       map[string]string
	query         map[string]string
	body          map[string]interface{}
	beforeHandler func(oRequest *Request) bool
	afterHandler  func(oRequest *Request, oResponseSet *ResponseSet)
}

func (oRequest Request) Name() string {
	return oRequest.name
}

func (oRequest *Request) SetName(sName string) {
	oRequest.name = sName
}

func (oRequest Request) Method() Method {
	return oRequest.method
}

func (oRequest *Request) SetMethod(oMethod Method) {
	oRequest.method = oMethod
}

func (oRequest Request) Endpoint() string {
	return oRequest.endpoint
}

func (oRequest *Request) SetEndpoint(sEndpoint string) {
	oRequest.endpoint = sEndpoint
}

func (oRequest Request) Headers() map[string]string {
	return oRequest.headers
}

func (oRequest *Request) SetHeaders(aHeaders map[string]string) {
	oRequest.headers = array.CloneMap(aHeaders)
}

func (oRequest *Request) AddHeader(sKey string, sValue string) {
	oRequest.headers[sKey] = sValue
}

func (oRequest Request) GetHeader(sKey string) string {
	return oRequest.headers[sKey]
}

func (oRequest *Request) DelHeader(sKey string) {
	delete(oRequest.headers, sKey)
}

func (oRequest Request) Query() map[string]string {
	return oRequest.query
}

func (oRequest *Request) SetQuery(aQuery map[string]string) {
	oRequest.query = array.CloneMap(aQuery)
}

func (oRequest *Request) AddQuery(sKey string, sValue string) {
	oRequest.query[sKey] = sValue
}

func (oRequest Request) GetQuery(sKey string) string {
	return oRequest.query[sKey]
}

func (oRequest *Request) DelQuery(sKey string) {
	delete(oRequest.query, sKey)
}

func (oRequest Request) Body() map[string]interface{} {
	return oRequest.body
}

func (oRequest *Request) SetBody(oBody map[string]interface{}) {
	oRequest.body = make(map[string]interface{})
	for sKey, oValue := range oBody {
		oRequest.body[sKey] = oValue
	}
}

func (oRequest *Request) AddBody(sKey string, oValue interface{}) {
	oRequest.body[sKey] = oValue
}

func (oRequest Request) GetBody(sKey string) interface{} {
	return oRequest.body[sKey]
}

func (oRequest *Request) DelBody(sKey string) {
	delete(oRequest.body, sKey)
}

func (oRequest *Request) SetBeforeHandler(fHandler func(oRequest *Request) bool) {
	oRequest.beforeHandler = fHandler
}

func (oRequest *Request) SetAfterHandler(fHandler func(oRequest *Request, oResponseSet *ResponseSet)) {
	oRequest.afterHandler = fHandler
}
