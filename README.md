# go-ecoflow-http

Go HTTP client library that implements authorization/signature procedure of the Ecoflow open platform API as specified in the dev documentation: https://developer-eu.ecoflow.com/us/document/generalInfo.
Wraps the standard library HTTP client.

## Usage

Make sure you have already obtained your developer access and secret keys. Then initialize the client:

```go

import efhttp "github.com/Drofff/go-ecoflow-http/pkg/http"

...

efConf := efhttp.ClientConfig{
    Host:      "https://api-e.ecoflow.com",
    AccessKey: "Fp4SvIprYSDPXtYJidEt*****",
    SecretKey: "WIbFEKre0s6sLnh4ei7SPUeYnp*****", // make sure to load secretly ;)
}
c := efhttp.NewClient(efConf, http.DefaultClient)

/*
For each individual request do the following:
 */

req, err := c.NewRequest(http.MethodGet, "/iot-open/sign/device/list", nil)
if err != nil {
	// handle error.
}

resp, err := c.Do(req)
// do your business with the `resp`.

```
