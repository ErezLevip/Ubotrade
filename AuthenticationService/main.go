package main

import (
	"net/http"
	"os"
	"context"
	"fmt"
	"strconv"
	"time"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Service"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Logging"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Instrumenting"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Transport"
	"github.com/erezlevip/Ubotrade/AuthenticationService/DependencyResolving"
	"github.com/erezlevip/Ubotrade/Global/ServiceHealth"
	"github.com/erezlevip/Ubotrade/Logger"
)


func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ubot",
		Subsystem: "authentication_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "authentication_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "authentication_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	machineName, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
	}

	serviceInfo := Global.ServiceInformation{
		Port:            8088,
		ServiceName:     "Authentication",
		Url:             "http://localhost",
		MachineName:     machineName,
		LastHealthCheck: time.Now(),
	}

	Logger.SetGlobalLogger()
	var svc AuthenticationService.IAuthenticationService
	svc = AuthenticationService.AuthenticationService{}.Init(serviceInfo)
	svc = DependencyResolving.Make(svc)
	svc = AuthenticationLogging.LoggingMiddleware{logger, svc}
	svc = AuthenticationInstrumenting.Make(requestCount, requestLatency, countResult, svc)

	ctx := DependencyResolving.GetRequestContext(context.Background())
	ServiceHealth.StartHealthTicker(ctx, serviceInfo)

	getLoginHandler := httptransport.NewServer(
		Transport.LoginEndpoint(svc),
		Transport.DecodeLoginRequest,
		Transport.EncodeResponse,
	)

	getTokenHandler := httptransport.NewServer(
		Transport.GetTokenEndpoint(svc),
		Transport.DecodeTokenRequest,
		Transport.EncodeResponse,
	)

	validateTokenHandler := httptransport.NewServer(
		Transport.ValidateTokenEndpoint(svc),
		Transport.DecodeTokenValidationRequest,
		Transport.EncodeResponse,
	)

	getRequestMetricsHandler := httptransport.NewServer(
		Transport.GetMetricsEndpoint(svc),
		Transport.DecodeMetricsRequest,
		Transport.EncodeResponse,
	)

	http.Handle("/login", getLoginHandler)
	http.Handle("/gettoken", getTokenHandler)
	http.Handle("/validatetoken", validateTokenHandler)
	http.Handle("/requestmetrics", getRequestMetricsHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":"+strconv.Itoa(serviceInfo.Port))
	logger.Log("err", http.ListenAndServe(":"+strconv.Itoa(serviceInfo.Port), nil))

}
