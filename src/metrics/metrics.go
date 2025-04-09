package metrics

func Initialize() {
	InitializeMemoryMetricies()
	InitializeUserMetricies()
	InitializeMessageMetricies()
	InitializeFollowerMetricies()
}