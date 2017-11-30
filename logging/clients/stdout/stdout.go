package stdout

import (
	"fmt"

	"github.com/turbosonic/api-gateway/logging"
)

type stdOutLogger struct {
}

func New() stdOutLogger {
	return stdOutLogger{}
}

func (std stdOutLogger) LogRequest(l *logging.RequestLog, index string, logType string) {
	fmt.Printf("%+v\n", *l)
}

func (std stdOutLogger) LogRelay(l *logging.RelayLog, index string, logType string) {
	fmt.Printf("%+v\n", *l)
}
