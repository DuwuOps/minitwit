package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeFollowerMetricies() {
	metrics := []prometheus.Collector{
		FollowerTotal,
		FolloweesTotal,
	}

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			log.Printf("❌ Error: Unable to register prometheus metric %T: %v", metric, err)
		}
	}
}

// Snapshots
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