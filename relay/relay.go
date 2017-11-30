package relay

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/turbosonic/api-gateway/logging"
)

type Relay struct {
	client *http.Client
	logger logging.LogClient
}

func New(MaxIdleConns int, IdleConnTimeout int, logger logging.LogClient) Relay {
	relay := Relay{}
	tr := &http.Transport{
		MaxIdleConns:       MaxIdleConns,
		IdleConnTimeout:    time.Duration(IdleConnTimeout) * time.Second,
		DisableCompression: true,
	}
	relay.client = &http.Client{Transport: tr}

	relay.logger = logger

	return relay
}

func (relay Relay) MakeRequest(r RelayRequest) (resp *http.Response, err error) {
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	start := time.Now()

	req.Header = r.Header

	resp, err = relay.client.Do(req)
	if err != nil {
		return nil, err
	}

	go func() {
		rl := logging.RelayLog{
			start,
			r.Header.Get("request_id"),
			r.Host,
			strings.Replace(r.URL, r.Host, "", 1),
			r.Method,
			resp.StatusCode,
			float64(time.Since(start)) / float64(time.Millisecond)}

		relay.logger.LogRelay(&rl, "api-gateway-request", "relay-request")
	}()

	return resp, nil
}

type RelayRequest struct {
	URL    string
	Method string
	Header http.Header
	Body   []byte
	Host   string
}
