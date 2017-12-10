package elk

import (
	"os"

	"github.com/turbosonic/api-gateway/logging"
	elastic "gopkg.in/olivere/elastic.v2"
)

type elkLogger struct {
	client *elastic.Client
}

func New() elkLogger {
	elasticURL := os.Getenv("LOGGING_ELASTIC_URL")
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticURL))
	if err != nil {
		panic(err)
	}

	return elkLogger{client}
}

func (elk elkLogger) LogRequest(l *logging.RequestLog, index string, logType string) {
	elk.client.Index().Index(index).Type(logType).Id(l.RequestID).BodyJson(l).Do()
}

func (elk elkLogger) LogRelay(l *logging.RelayLog, index string, logType string) {
	elk.client.Index().Index(index).Type(logType).BodyJson(l).Do()
}
