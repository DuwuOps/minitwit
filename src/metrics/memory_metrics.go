package metrics

import (
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/mem"
)

func InitializeMemoryMetricies() {
	if err := prometheus.Register(MemoryUsage); err != nil {
		utils.LogError("Unable to register prometheus metric MemoryUsage", err)
	}
}

var (
	MemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "minitwit_cpu_usage_gauge",
			Help: "CPU Usage",
		},
		[]string{"parameter"},
	)
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {			
			err := next(c)

			vm, _ := mem.VirtualMemory()
			MemoryUsage.WithLabelValues("UsedPercent").Set(vm.UsedPercent)
			MemoryUsage.WithLabelValues("Used").Set(float64(vm.Used))
			MemoryUsage.WithLabelValues("Available").Set(float64(vm.Available))
			MemoryUsage.WithLabelValues("Total").Set(float64(vm.Total))
			
			return err
		}
	}
}