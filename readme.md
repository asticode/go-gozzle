# About

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/asticode/go-gozzle/gozzle)

`go-gozzle` is an HTTP client wrapper that ease sending multiple requests asynchronously and processing their respective responses for the GO programming language (http://golang.org).

# Install `go-gozzle`

Run the following command:

    $ go get github.com/asticode/go-gozzle/gozzle
    
# Example

    import (
        "github.com/asticode/go-gozzle/gozzle"
    )
    
    // Create gozzle without a maximum body size for the response
    g := gozzle.NewGozzle(0)
    
    // Create a request set
    reqSet := gozzle.NewRequestSet()
    
    // Create the first request
    r := NewRequest("my-first-request", gozzle.MethodGet, "/my-first-endpoint")
    
    // Add headers
    r.AddHeader("X-MyHeader", "MyValue")
    
    // Add query
    r.SetQuery(map[string]string{
        "api_key": "my_api_key",
        "access_token": "my_access_token",
    })
    
    // Set the callback that is executed before sending the request
    r.SetBeforeHandler(func(r gozzle.Request) bool {
        if r.GetHeader("X-MyHeader") != "MyValue" {
            // Request will be sent
            return true
        } else {
            // Request will not be sent
            return false
        }
    })
    
    // Add the request to the request set
    reqSet.AddRequest(r)
    
    // Create the second request
    r = NewRequest("my-second-request", gozzle.MethodPost, "/my-second-endpoint")
    
    // Add a body to the request
    r.SetBody(map[string]interface{}{
        "name": "Asticode",
        "email": "test@asticode.com",
    })
    
    // Set the callback that is executed after sending the request
    r.SetAfterHandler(func(req gozzle.Request, resp gozzle.Response) {
        b, e := ioutil.ReadAll(resp.BodyReader())
        fmt.Println(fmt.Sprintf("Body %s for path %s", string(b), req.Path())
    })
    
    // Add the request to the request set
    reqSet.AddRequest(r)
    
    // Execute the 2 requests asynchronously
    respSet := g.Exec(reqSet)
    
    // Process responses
    // The first request has not been sent because of the beforeHandler so only one response should be stored
    for _, name := range respSet.Names() {
        req := reqSet.GetRequest(name)
        resp := respSet.GetResponse(name)
        if len(resp.Errors()) > 0 {
            fmt.Println(fmt.Sprintf("Error in request to %s", req.Path()))
        } else {
            fmt.Println(fmt.Sprintf("Request to %s was successful", req.Path()))
        }
    }
    