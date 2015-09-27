package gozzle

type RequestSet map[string]*Request

func (oRequestSet *RequestSet) AddRequest(sName string, oMethod Method, sEndpoint string) {
	(*oRequestSet)[sName] = &Request{
		name:     sName,
		method:   oMethod,
		endpoint: sEndpoint,
	}
}
