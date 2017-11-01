package main

import (
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"local/UbotTrade/DataHandlers/Redis"
	"local/UbotTrade/RegistryService/Instrumenting"
	"local/UbotTrade/RegistryService/Logging"
	"local/UbotTrade/RegistryService/Service"
	"local/UbotTrade/RegistryService/Transport"
	"net/http"
	"os"
	"local/UbotTrade/RegistryService/DependencyResolving"
	"local/UbotTrade/Logger"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ubot",
		Subsystem: "registry_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "registry_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "registry_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	var svc RegistryService.IRegistryService
	Logger.SetGlobalLogger()
	svc = RegistryService.RegistryService{}
	svc = DependencyResolving.Make(svc)
	svc = RegistryServiceLogging.LoggingMiddleware{logger, svc}
	svc = RegistryServiceInstrumenting.InstrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	//svc = DependencyResolving.Make(svc)

	svc.Init(DataHandlers.RedisConfiguration{
		ConnectionString: "localhost:6379",
		Credentials:      "",
		Db:               1,
	})
	registerHandler := httptransport.NewServer(
		RegistryServiceTransport.RegisterEndpoint(svc),
		RegistryServiceTransport.DecodeRegistryRequest,
		RegistryServiceTransport.EncodeResponse,
	)

	deRegisterHandler := httptransport.NewServer(
		RegistryServiceTransport.DeRegisterEndpoint(svc),
		RegistryServiceTransport.DecodeRegistryRequest,
		RegistryServiceTransport.EncodeResponse,
	)

	getServiceHandler := httptransport.NewServer(
		RegistryServiceTransport.GetServiceEndpoint(svc),
		RegistryServiceTransport.DecodeRegistryRequest,
		RegistryServiceTransport.EncodeResponse,
	)

	http.Handle("/register", registerHandler)
	http.Handle("/deregister", deRegisterHandler)
	http.Handle("/getservice", getServiceHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
