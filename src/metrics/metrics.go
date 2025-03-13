package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalCallCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "minitwit_route_calls_total",
			Help: "Total route calls",
		},
		[]string{"route", "http_method", "status_code", "duration"},
	)
)

func Initialize() error {
	if err := prometheus.Register(TotalCallCounter); err != nil {
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
			
			return err
		}
	}
}