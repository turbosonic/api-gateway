package elk

import (
	"github.com/turbosonic/api-gateway/logging"
	elastic "gopkg.in/olivere/elastic.v2"
)

type elkLogger struct {
	client *elastic.Client
}

func New(elasticURL string) elkLogger {
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticURL))
	if err != nil {
		panic(err)
	}

	return elkLogger{client}
}

func (elk elkLogger) Log(l *logging.Log, index string, logType string) {
	elk.client.Index().Index(index).Type(logType).Id(l.RequestID).BodyJson(l).Do()
}
