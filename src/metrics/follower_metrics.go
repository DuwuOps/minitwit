package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeFollowerMetricies() {
	metrics := []prometheus.Collector{
		FollowerTotal,
		FolloweesTotal,
	}

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			slog.Error("Unable to register prometheus metric", slog.Any("error", err), slog.Any("metric", metric))
		}
	}
}

// Snapshots.
var FollowerTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "minitwit_follower_total",
		Help: "Number of users with a specified follower count range",
	},
	[]string{"range"},
)

var FolloweesTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "minitwit_followees_total",
		Help: "Number of users with a specified followee count range",
	},
	[]string{"range"},
)
