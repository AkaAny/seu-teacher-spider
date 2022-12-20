package main

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"net/http"
)

const UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"

type UARoundTripper struct {
	base http.RoundTripper
	UA   string
}

func (U *UARoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("User-Agent", U.UA)
	return U.base.RoundTrip(request)
}

func NewUARoundTripper(base http.RoundTripper, ua string) *UARoundTripper {
	return &UARoundTripper{
		base: base,
		UA:   ua,
	}
}

func NewGoutClient() *dataflow.Gout {
	var httpClient = &http.Client{}
	httpClient.Transport = NewUARoundTripper(http.DefaultTransport, UA)
	return gout.New(httpClient)
}
