package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeMessageMetricies() {
	metrics := []prometheus.Collector{
		MessagesPosts,
	}

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			log.Printf("❌ Error: Unable to register prometheus metric %T: %v", metric, err)
		}
	}
}

var MessagesPosts = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "minitwit_messages_counter",
		Help: "Number of messages posted, labeled by hour and weekday.",
	},
	[]string{"hour", "weekday"},
)
