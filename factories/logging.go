package factories

import (
	"fmt"
	"os"

	"github.com/turbosonic/api-gateway/logging"
	"github.com/turbosonic/api-gateway/logging/clients/elk"
	"github.com/turbosonic/api-gateway/logging/clients/stdout"
)

func LogClient() logging.LogClient {
	authProvider := os.Getenv("LOGGING_PROVIDER")

	switch authProvider {
	case "elk":
		fmt.Println("[x] Logging to ELK")
		return elk.New()
	default:
		fmt.Println("[x] Logging to stdout")
		return stdout.New()
	}
}
