package factories

import (
	"fmt"
	"os"

	"github.com/turbosonic/api-gateway/logging"
	"github.com/turbosonic/api-gateway/logging/clients/elasticsearch"
	"github.com/turbosonic/api-gateway/logging/clients/stdout"
)

func LogClient() logging.LogClient {
	authProvider := os.Getenv("LOGGING_PROVIDER")

	switch authProvider {
	case "elasticsearch":
		fmt.Println("[x] Logging to Elasticsearch")
		return elasticsearch.New()
	default:
		fmt.Println("[x] Logging to stdout")
		return stdout.New()
	}
}
