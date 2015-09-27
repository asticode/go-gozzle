# About

This is a HTTP client library for the GO programming language (http://golang.org).

For the parallelizator see https://github.com/asticode/go-parallelizator

# Dependencies

    github.com/asticode/go-parallelizator/parallelizator
    github.com/asticode/go-toolbox/array

# Installing

## Using *go get*

    $ go get github.com/asticode/go-gozzle/gozzle
    
After this command *go-gozzle* is ready to use. Its source will be in:

    $GOPATH/src/github.com/asticode/go-gozzle/gozzle
    
# Configuration

Best is to use JSON configuration:

    {
        "base_url": "BASE_URL"
    }
    
And decode it:

    oGozzleConfiguration = json.UnMarshall(sConfiguration)
    
An example of the configuration would be:

    {
        "base_url": "http://mysite.com"
    }
    
# Example

    import (
        "github.com/asticode/go-gozzle/gozzle"
    )

    // Create Gozzle
    oGozzle, oErr := gozzle.NewGozzle(oGozzleConfiguration, oParallelizator)
    
    // Initialize request set
    oRequestSet := gozzle.RequestSet{}
    
    // Add the first request
    oRequestSet.AddRequest("my-first-request", gozzle.METHOD_GET, "/my-first-endpoint")
    
    // Set headers of the first request
    oRequestSet["my-first-request"].AddHeader("X-MyHeader", "MyValue")
    
    // Set query of the first request
    oRequestSet["my-first-request"].SetQuery(map[string]string{
        "api_key": "my_api_key",
        "access_token": "my_access_token",
    })
    
    // Set callback executed before sending the first request
    oRequestSet["my-first-request"].SetBeforeHandler(func(oRequest *gozzle.Request) bool {
        if MyVar == "MyValue" {
            // Request will be sent
            return true
        } else {
            // Request will not be sent
            return false
        }
    })
    
    // Add the second request
    oRequestSet.AddRequest("my-second-request", gozzle.METHOD_POST, "/my-second-endpoint")
    
    // Set body of the second request
    oRequestSet["my-second-request"].SetBody(map[string]interface{}{
        "name": "Asticode",
        "email": "test@asticode.com",
    })
    
    // Set callback executed after sending the second request
    oRequestSet["my-second-request"].SetAfterHandler(func(oRequest *gozzle.Request, oResponseSet *gozzle.ResponseSet) bool {
        fmt.Println((*(*oResponseSet)[oRequest.Name()].Body()))
    })
    
    // Exec request set
    oResponseSet, oErr := oGozzle.Exec(oRequestSet)
    
    // Process responses
    for sResponseName, oResponse := range *oResponseSet {
        if oResponse.IsError() {
            fmt.Println(fmt.Sprintf("Error in %s request", sResponseName))
        } else {
            fmt.Println(fmt.Sprintf("Request %s was successful", sResponseName))
        }
    }
    