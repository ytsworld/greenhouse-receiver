package receiver

import (

	// Imports the Stackdriver Monitoring client package.
	"fmt"
	"log"
	"os"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	greenhouse "github.com/ytsworld/greenhouse-client/pkg"
	"golang.org/x/net/context"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

var (
	client    *monitoring.MetricClient
	projectID string
)

func init() {
	var err error

	client, err = monitoring.NewMetricClient(context.Background())
	if err != nil {
		log.Fatalf("monitoring.NewMetricClient: %v", err)
	}

	// When running in a cloud function environment the environment variable is always set
	projectID = os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.Fatalf("missing project id in env var GCP_PROJECT")
	}
}

func persistAll(data greenhouse.Data) error {
	//TODO which timestamp to use? Date on raspi could be out of sync?
	timestamp := data.UnixTimestampUTC

	series := []*monitoringpb.TimeSeries{
		createTimeSeriesForData("temperature", timestamp, float64(data.Temperature)),
		createTimeSeriesForData("humidity", timestamp, float64(data.Humidity)),
		// TODO moisture data needs interpretation
		createTimeSeriesForData("soil_moisture", timestamp, float64(data.SoilMoistureResistance)),
	}

	err := persistTimeSeries(series)
	if err != nil {
		return fmt.Errorf("error persisting data: %v", err)
	}

	return nil
}

func persistTimeSeries(timeSeries []*monitoringpb.TimeSeries) error {
	log.Printf("Project id: %s Name: %s", projectID, monitoring.MetricProjectPath(projectID))
	// Writes time series data.
	if err := client.CreateTimeSeries(context.Background(), &monitoringpb.CreateTimeSeriesRequest{
		Name:       monitoring.MetricProjectPath(projectID),
		TimeSeries: timeSeries,
	}); err != nil {
		log.Fatalf("Failed to write time series data: %v", err)
	}

	fmt.Printf("Done writing time series data.\n")

	return nil
}

func createTimeSeriesForData(metricName string, timestamp int64, dataPoint float64) *monitoringpb.TimeSeries {
	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type: fmt.Sprintf("custom.googleapis.com/greenhouse/%s", metricName),
			Labels: map[string]string{
				"instance": "one",
			},
		},
		Resource: &monitoredrespb.MonitoredResource{
			Type: "global",
			Labels: map[string]string{
				"project_id": projectID,
			},
		},
		Points: []*monitoringpb.Point{
			&monitoringpb.Point{
				Interval: &monitoringpb.TimeInterval{
					EndTime: &googlepb.Timestamp{
						Seconds: timestamp,
					},
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_DoubleValue{
						DoubleValue: float64(dataPoint),
					},
				},
			},
		},
	}
}
