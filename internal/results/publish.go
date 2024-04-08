package results

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const metricName = "latency"

type Publish struct {
	pushGateway      string
	cloudWatchPrefix string
	cloudWatchClient *cloudwatch.Client
}

func NewPublish(pushGateway string, cloudWatchPrefix string, cloudWatchClient *cloudwatch.Client) *Publish {
	return &Publish{
		pushGateway:      pushGateway,
		cloudWatchPrefix: cloudWatchPrefix,
		cloudWatchClient: cloudWatchClient,
	}
}

func (p *Publish) Send(metrics IResults) error {
	// prometheus pushgateway
	if p.pushGateway != "" {
		slog.Debug("Sending metrics to Pushgateway", "gateway", p.pushGateway)
		pusher := push.New(p.pushGateway, metrics.GetId())
		ctx, cncl := context.WithTimeout(context.Background(), time.Second*3)
		defer cncl()
		err := pusher.Collector(metrics.GetPrometheusMetrics()).PushContext(ctx)
		if err != nil {
			slog.Error("Could not push completion time to Pushgateway:", "error", err)
			return err
		}
	}
	// cloudwatch
	if p.cloudWatchPrefix != "" {
		slog.Debug("Sending metrics to CloudWatch", "namespace", p.cloudWatchPrefix+"/"+metrics.GetId())
		_, err := p.cloudWatchClient.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(p.cloudWatchPrefix + "/" + metrics.GetId()),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(metricName),
					Value:      aws.Float64(metrics.GetLatency()),
					Dimensions: metrics.GetCloudWatchDimensions(),
				},
			},
		})
		if err != nil {
			slog.Error("Could not send metric data to CloudWatch:", "error", err)
			return err
		}
	}
	return nil
}
