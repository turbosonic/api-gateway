package applicationinsights

import (
	"os"
	"strconv"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/turbosonic/api-gateway/logging"
)

type applicationinsightsLogger struct {
	client appinsights.TelemetryClient
}

func New() applicationinsightsLogger {
	aiInsKey := os.Getenv("APPLICATIONINSIGHTS_INTRUMENTATION_KEY")
	client := appinsights.NewTelemetryClient(aiInsKey)

	return applicationinsightsLogger{client}
}

func (ai applicationinsightsLogger) LogRequest(l *logging.RequestLog, index string, logType string) {
	duration := time.Duration(l.Duration * 1000000)
	request := appinsights.NewRequestTelemetry(l.Method, l.URL, duration, strconv.Itoa(l.StatusCode))
	request.Id = l.RequestID
	if l.StatusCode > 399 {
		request.Success = false
	} else {
		request.Success = true
	}
	request.Properties["Config"] = l.Config
	request.Timestamp = l.Date
	request.Properties["Path"] = l.Path
	request.Properties["Agent"] = l.Agent
	request.Properties["Host"] = l.Host
	request.Properties["OS"] = l.OS
	request.Properties["RemoteAddr"] = l.RemoteAddr
	request.Properties["GoVersion"] = l.GoVersion

	request.Measurements["Content-Length"] = float64(l.ContentLength)

	ai.client.Track(request)
}

func (ai applicationinsightsLogger) LogRelay(l *logging.RelayLog, index string, logType string) {
	success := true
	if l.StatusCode > 399 {
		success = false
	}

	dependency := appinsights.NewRemoteDependencyTelemetry(l.Host, l.Method, l.URL, success)
	dependency.Id = l.RequestID
	dependency.ResultCode = strconv.Itoa(l.StatusCode)
	dependency.Duration = time.Duration(l.Duration * 1000000)
	dependency.Timestamp = l.Date

	ai.client.Track(dependency)
}
