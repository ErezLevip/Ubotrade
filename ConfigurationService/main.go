package main

import (
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"

	"fmt"
	"github.com/go-kit/kit/log"
	"local/UbotTrade/ConfigurationService/Instrumenting"
	"local/UbotTrade/ConfigurationService/Logging"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/ConfigurationService/Transport"
	"local/UbotTrade/Global"
	"strconv"
	"time"
	"local/UbotTrade/Global/ServiceHealth"
	"local/UbotTrade/ConfigurationService/DependencyResolving"
	"context"
	"local/UbotTrade/Logger"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ubot",
		Subsystem: "configuration_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "configuration_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "configuration_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	machineName, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
	}
	serviceInfo := Global.ServiceInformation{Port: 8081,
		ServiceName:     "Configuration",
		Url:             "http://localhost",
		MachineName:     machineName,
		LastHealthCheck: time.Now(),
	}

	Logger.SetGlobalLogger()

	staticConfigPath := "/ConfigurationService/BasicConfiguration.json"
	var svc ConfigurationService.IConfigurationService
	svc = ConfigurationService.ConfigurationService{}.Init(staticConfigPath, serviceInfo)
	svc = DependencyResolving.Make(svc,staticConfigPath)
	svc = ConfigurationServiceLogging.LoggingMiddleware{logger, svc}
	svc = *ConfigurationServiceInstrumenting.Make(requestCount, requestLatency, countResult, svc)


	ctx := DependencyResolving.GetRequestContext(context.Background(),staticConfigPath)
	ServiceHealth.StartHealthTicker(ctx, serviceInfo)

	getConfigHandler := httptransport.NewServer(
		ConfigurationServiceTransport.GetConfigurationEndpoint(svc),
		ConfigurationServiceTransport.DecodeConfigurationRequest,
		ConfigurationServiceTransport.EncodeResponse,
	)
	getRequestMetricsHandler := httptransport.NewServer(
		ConfigurationServiceTransport.GetMetricsEndpoint(svc),
		ConfigurationServiceTransport.DecodeMetricsRequest,
		ConfigurationServiceTransport.EncodeResponse,
	)

	http.Handle("/getconfiguration", getConfigHandler)
	http.Handle("/requestmetrics", getRequestMetricsHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":", serviceInfo.Port)
	logger.Log("err", http.ListenAndServe(":"+strconv.Itoa(serviceInfo.Port), nil))
}
