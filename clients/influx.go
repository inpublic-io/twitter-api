package clients

import (
	"log"
	"os"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

var (
	influxClient influxdb.Client
	influxOrg    string
	influxBucket string
)

func InitializeInfluxClient() {
	influxToken, ok := os.LookupEnv("INFLUX_TOKEN")
	if !ok {
		log.Fatalf("INFLUX_TOKEN env var is required")
	}

	influxOrg, ok = os.LookupEnv("INFLUX_ORG")
	if !ok {
		log.Fatalf("INFLUX_ORG env var is required")
	}

	influxBucket, ok = os.LookupEnv("INFLUX_BUCKET")
	if !ok {
		log.Fatalf("INFLUX_BUCKET env var is required")
	}

	influxClient = influxdb.NewClient("https://us-central1-1.gcp.cloud2.influxdata.com", influxToken)
}

func InfluxClient() influxdb.Client {
	return influxClient
}

func InfluxOrg() string {
	return influxOrg
}

func InfluxBucket() string {
	return influxBucket
}
