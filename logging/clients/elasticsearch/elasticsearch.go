package elasticsearch

import (
	"os"

	"github.com/turbosonic/api-gateway/logging"
	elastic "gopkg.in/olivere/elastic.v2"
)

type elasticsearchLogger struct {
	client *elastic.Client
}

func New() elasticsearchLogger {
	elasticURL := os.Getenv("LOGGING_ELASTIC_URL")
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticURL))
	if err != nil {
		panic(err)
	}

	return elasticsearchLogger{client}
}

func (elasticsearch elasticsearchLogger) LogRequest(l *logging.RequestLog, index string, logType string) {
	elasticsearch.client.Index().Index(index).Type(logType).Id(l.RequestID).BodyJson(l).Do()
}

func (elasticsearch elasticsearchLogger) LogRelay(l *logging.RelayLog, index string, logType string) {
	elasticsearch.client.Index().Index(index).Type(logType).BodyJson(l).Do()
}
