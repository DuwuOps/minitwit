package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/mem"
)

var (
	TotalCallCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "minitwit_route_calls_total",
			Help: "Total route calls",
		},
		[]string{"route", "http_method", "status_code", "duration"},
	)

	CpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "minitwit_cpu_usage_gauge",
			Help: "CPU Usage",

		},
		[]string{"parameter"},
	)


)

func Initialize() error {
	if err := prometheus.Register(TotalCallCounter); err != nil {
		return err
	}

	if err := prometheus.Register(CpuUsage); err != nil {
		return err
	}
	
	return nil
}

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			route := c.Path()
			requestHttpMethod := c.Request().Method
			responseStatusCode := c.Response().Status
			request_duration := time.Since(start).Milliseconds()
			
			// Increment the request counter
			TotalCallCounter.WithLabelValues(
				route,
				requestHttpMethod,
				strconv.Itoa(responseStatusCode),
				strconv.FormatInt(request_duration, 10),
			).Inc()


			v, _ := mem.VirtualMemory()
			fmt.Printf("v.UsedPercent: %v\n", v.UsedPercent)
			CpuUsage.WithLabelValues("UsedPercent").Set(v.UsedPercent)
			CpuUsage.WithLabelValues("Used").Set(float64(v.Used))
			CpuUsage.WithLabelValues("Available").Set(float64(v.Available))
			CpuUsage.WithLabelValues("Total").Set(float64(v.Total))
			
			return err
		}
	}
}