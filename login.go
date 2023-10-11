/*
thehive5 implements functionality to interact with the most recent version of thehive.
https://www.strangebee.com/thehive/
*/
package thehive5

import (
	"crypto/tls"
	"net/http"
	"strings"
)

// A Hivedata stores the apikey, url and http client for subsequent API calls
type Hivedata struct {
	Url    string
	Apikey string
	Client HttpClient
}

// a HttpClient interface gets used for testing
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CreateLogin
// Defines API login principles that can be reused in requests
// Returns a Hivedata struct
func CreateLogin(inurl string, apikey string, verify bool) Hivedata {
	return Hivedata{
		Url:    strings.TrimRight(inurl, "/"),
		Apikey: apikey,
		Client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: !verify},
			},
		},
	}
}
