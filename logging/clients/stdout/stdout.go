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

func (std stdOutLogger) Log(l *logging.Log, index string, logType string) {
	fmt.Printf("%+v\n", *l)
}
