package grpc

import (
	"context"
	"fmt"
	"time"

	twitterv1 "github.com/inpublic-io/inpublicapis/twitter/v1"
	"github.com/inpublic-io/twitter-api/clients"
)

// metricsServer is used to implement io.inpublic.twitter.v1.Metrics
type metricsServer struct {
	twitterv1.UnimplementedMetricsServer
}

func (s *metricsServer) ContributionsPerInterval(ctx context.Context, req *twitterv1.IntervalAggregateQuery) (*twitterv1.ContributionMetric, error) {
	influxClient := clients.InfluxClient()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := influxClient.QueryAPI(clients.InfluxOrg()).Query(ctx,
		fmt.Sprintf(`from(bucket: "hashtag-metrics")
			|> range(start: %s)
			|> group(columns: ["_time"])
			|> aggregateWindow(every: %s, fn: mean, createEmpty: false)
			|> count()
	`, req.Range, req.Interval))
	if err != nil {
		return nil, err
	}

	var metrics = make(map[string]int32)

	// Iterate over query response
	for result.Next() {
		metrics[result.Record().Time().Format(time.RFC3339)] = int32(result.Record().Value().(int64))
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %v\n", result.Err().Error())
	}

	return &twitterv1.ContributionMetric{
		Values: metrics,
	}, nil
}
