package influxdb

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/turbosonic/api-gateway/logging"
)

var (
	influxDBNAME string
	wg           sync.WaitGroup
	bp           client.BatchPoints
)

type influxdbLogger struct {
	client client.Client
}

func New() influxdbLogger {
	influxDBURL := os.Getenv("LOGGING_INFLUXDB_URL")
	if influxDBURL == "" {
		panic("No LOGGING_INFLUXDB_URL environment variable found")
	}

	influxDBNAME = os.Getenv("LOGGING_INFLUX_DB_NAME")
	if influxDBNAME == "" {
		panic("No LOGGING_INFLUX_DB_NAME environment variable found")
	}

	// create a client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxDBURL,
	})
	if err != nil {
		panic("Error creating InfluxDB Client")
	}

	// need to create the database here, unless it already exists
	q := client.NewQuery("CREATE DATABASE "+influxDBNAME, "", "")
	if response, err := c.Query(q); err != nil && response.Error() != nil {
		fmt.Println(response.Error())
	}

	bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDBNAME,
		Precision: "s",
	})

	logger := influxdbLogger{
		client: c,
	}

	logger.startDataWriter()

	return logger
}

func (influxdb influxdbLogger) LogRequest(l *logging.RequestLog, index string, logType string) {
	tags := map[string]string{
		"RequestID":     l.RequestID,
		"Config":        l.Config,
		"Path":          l.Path,
		"URL":           l.URL,
		"Method":        l.Method,
		"StatusCode":    strconv.FormatInt(int64(l.StatusCode), 10),
		"ContentLength": strconv.FormatInt(int64(l.ContentLength), 10),
		"Host":          l.Host,
		"RemoteAddr":    l.RemoteAddr,
		"Agent":         l.Agent,
		"OS":            l.OS,
		"GoVersion":     l.GoVersion,
	}
	fields := map[string]interface{}{
		"duration": l.Duration,
	}
	pt, err := client.NewPoint("request", tags, fields, l.Date)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	wg.Wait()
	bp.AddPoint(pt)
}

func (influxdb influxdbLogger) LogRelay(l *logging.RelayLog, index string, logType string) {
	tags := map[string]string{
		"RequestID":  l.RequestID,
		"URL":        l.URL,
		"Host":       l.Host,
		"Method":     l.Method,
		"StatusCode": strconv.FormatInt(int64(l.StatusCode), 10),
	}
	fields := map[string]interface{}{
		"duration": l.Duration,
	}
	pt, err := client.NewPoint("relay", tags, fields, l.Date)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	wg.Wait()
	bp.AddPoint(pt)
}

func (influxdb influxdbLogger) startDataWriter() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			wg.Add(1)
			pointsCount := len(bp.Points())
			if pointsCount > 0 {
				influxdb.client.Write(bp)
				bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
					Database:  influxDBNAME,
					Precision: "s",
				})

				fmt.Printf("%d points sent to InfluxDB\n", pointsCount)
			}
			wg.Done()
		}
	}()
}
