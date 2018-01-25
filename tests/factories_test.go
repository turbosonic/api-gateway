package factories

import (
	"testing"

	"github.com/turbosonic/api-gateway/factories"
)

func TestGetLogger(t *testing.T) {
	var logger = factories.LogClient()

	if logger == nil {
		t.Errorf("Logger could not be created")
	}
}
