package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeUserMetricies() {
	metrics := []prometheus.Collector{
		NewUsers,
		TotalUsers,
	}

	for _, metric := range metrics{
		if err := prometheus.Register(metric); err != nil {
			slog.Error("Unable to register prometheus metric", slog.Any("error", err), slog.Any("metric", metric))
		}
	}
}
var NewUsers = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "minitwit_new_users_created",
        Help: "Number of new users created.",
    },
)

var TotalUsers = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Name: "minitwit_users_created_total",
        Help: "Total number of current users created.",
    },
)
