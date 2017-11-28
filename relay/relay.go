package relay

import (
	"bytes"
	"net/http"
	"time"
)

type Relay struct {
	client *http.Client
}

func New(MaxIdleConns int, IdleConnTimeout int) Relay {
	relay := Relay{}
	tr := &http.Transport{
		MaxIdleConns:       MaxIdleConns,
		IdleConnTimeout:    time.Duration(IdleConnTimeout) * time.Second,
		DisableCompression: true,
	}
	relay.client = &http.Client{Transport: tr}

	return relay
}

func (relay Relay) MakeRequest(r RelayRequest) (resp *http.Response, err error) {
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	req.Header = r.Header

	resp, err = relay.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type RelayRequest struct {
	URL    string
	Method string
	Header http.Header
	Body   []byte
}
