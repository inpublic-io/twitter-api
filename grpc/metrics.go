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

func (s *metricsServer) ContributionsPerInterval(ctx context.Context, req *twitterv1.IntervalAggregateQuery) (*twitterv1.ContributionMetrics, error) {
	influxClient := clients.InfluxClient()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := influxClient.QueryAPI(clients.InfluxOrg()).Query(ctx,
		fmt.Sprintf(`from(bucket: "%s")
			|> range(start: %s)
			|> group(columns: ["_time"])
			|> aggregateWindow(every: %s, fn: mean, createEmpty: false)
			|> count()
	`, clients.InfluxBucket(), req.Range, req.Interval))
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

	return &twitterv1.ContributionMetrics{
		Values: metrics,
	}, nil
}

func (s *metricsServer) ContributionsStats(ctx context.Context, req *twitterv1.RangeAggregateQuery) (*twitterv1.StatsMetrics, error) {
	influxClient := clients.InfluxClient()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := influxClient.QueryAPI(clients.InfluxOrg()).Query(ctx,
		fmt.Sprintf(`from(bucket: "%s")
		|> range(start: %s)
		|> group()
		|> aggregateWindow(every: 1d, fn: count)
		|> reduce(
		  identity: {
			total: 1,
			sum: 0,
			average: 0
		  },
		  fn: (r, accumulator) => ({
			total: accumulator.total + 1,
			sum:   accumulator.sum + r._value,
			average:   accumulator.sum / accumulator.total
		  })
		)
		|> drop(columns: ["total"])
	`, clients.InfluxBucket(), req.Range))
	if err != nil {
		return nil, err
	}

	metrics := &twitterv1.StatsMetrics{}

	for result.Next() {
		values := result.Record().Values()
		metrics.Average = int32(values["average"].(int64))
		// metrics.Total = int32(values["total"].(int64))
		metrics.Sum = int32(values["sum"].(int64))
	}
	if result.Err() != nil {
		fmt.Printf("query parsing error: %v\n", result.Err().Error())
	}

	return metrics, nil
}

func (s *metricsServer) ContributorsReach(ctx context.Context, req *twitterv1.RangeAggregateQuery) (*twitterv1.ReachMetrics, error) {
	influxClient := clients.InfluxClient()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := influxClient.QueryAPI(clients.InfluxOrg()).Query(ctx,
		fmt.Sprintf(`from(bucket: "%s")
		|> range(start: %s)
		|> group()
		|> keep(columns: ["_value", "author"])
		|> unique(column: "author")
		|> keep(columns: ["_value"])
	`, clients.InfluxBucket(), req.Range))
	if err != nil {
		return nil, err
	}

	var metrics []int32

	for result.Next() {
		metrics = append(metrics, int32(result.Record().Value().(int64)))
	}
	if result.Err() != nil {
		fmt.Printf("query parsing error: %v\n", result.Err().Error())
	}

	return &twitterv1.ReachMetrics{
		Values: metrics,
	}, nil
}

func (s *metricsServer) ContributorsStats(ctx context.Context, req *twitterv1.RangeAggregateQuery) (*twitterv1.StatsMetrics, error) {
	influxClient := clients.InfluxClient()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := influxClient.QueryAPI(clients.InfluxOrg()).Query(ctx,
		fmt.Sprintf(`from(bucket: "%s")
		|> range(start: %s)
		|> group()
		|> keep(columns: ["_value", "author"])
		|> unique(column: "author")
		|> keep(columns: ["_value"])
		|> reduce(
			identity: {
			  total: 1,
			  sum:   0,
			  average:   0
			},
			fn: (r, accumulator) => ({
			  total: accumulator.total + 1,
			  sum:   accumulator.sum + r._value,
			  average:   accumulator.sum / accumulator.total
			})
		)
	`, clients.InfluxBucket(), req.Range))
	if err != nil {
		return nil, err
	}

	metrics := &twitterv1.StatsMetrics{}

	for result.Next() {
		values := result.Record().Values()
		metrics.Average = int32(values["average"].(int64))
		metrics.Total = int32(values["total"].(int64))
		metrics.Sum = int32(values["sum"].(int64))
	}
	if result.Err() != nil {
		fmt.Printf("query parsing error: %v\n", result.Err().Error())
	}

	return metrics, nil
}
