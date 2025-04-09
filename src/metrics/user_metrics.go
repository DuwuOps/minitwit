package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeUserMetricies() {
	metrics := []prometheus.Collector{
		NewUsers,
	}

	for _, metric := range metrics{
		if err := prometheus.Register(metric); err != nil {
			log.Printf("❌ Error: Unable to register prometheus metric %T: %v", metric, err)
		}
	}
}
var NewUsers = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "minitwit_new_users_created",
        Help: "Number of new users created.",
    },
)

