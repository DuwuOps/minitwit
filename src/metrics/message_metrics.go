package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

func InitializeMessageMetricies() {
	metrics := []prometheus.Collector{
		MessagesPosts,
		MessagesTotal,
		FlaggedMessagesTotal,
	}

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			slog.Error("Unable to register prometheus metric", slog.Any("error", err), slog.Any("metric", metric))
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

// Snapshots.
var MessagesTotal = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "minitwit_messages_total",
		Help: "Number of messages posted, labeled by authorID.",
	},
)

var FlaggedMessagesTotal = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "minitwit_flagged_messages_total",
		Help: "Number of flagged messages, labeled by authorID.",
	},
)
