package main

import (
	"net/http"
	"os"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/UserService/Service"
	"github.com/erezlevip/Ubotrade/UserService/Logging"
	"github.com/erezlevip/Ubotrade/UserService/Transport"
	"github.com/erezlevip/Ubotrade/UserService/Instrumenting"
	"github.com/erezlevip/Ubotrade/Global/ServiceHealth"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/erezlevip/Ubotrade/UserService/DependencyResolving"
	"github.com/erezlevip/Ubotrade/Logger"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ubot",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "user_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "user_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	machineName, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
	}

	serviceInfo := Global.ServiceInformation{
		Port:            8091,
		ServiceName:     "User",
		Url:             "http://localhost",
		MachineName:     machineName,
		LastHealthCheck: time.Now(),
	}

	Logger.SetGlobalLogger()
	var svc UserService.IUserService
	svc = UserService.UserService{}.Init(serviceInfo)
	svc = DependencyResolving.Make(svc)
	svc = UserLogging.LoggingMiddleware{logger, svc}
	svc = UserInstrumenting.Make(requestCount, requestLatency, countResult, svc)

	ctx := DependencyResolving.GetRequestContext(context.Background())
	ServiceHealth.StartHealthTicker(ctx, serviceInfo)

	getCreateUserHandler := httptransport.NewServer(
		Transport.CreateUserEndpoint(svc),
		Transport.DecodeCreateUserRequest,
		Transport.EncodeResponse,
	)

	getGetUserHandler := httptransport.NewServer(
		Transport.GetUserEndpoint(svc),
		Transport.DecodeGetUserRequest,
		Transport.EncodeResponse,
	)

	getSetUserHandler := httptransport.NewServer(
		Transport.SetUserEndpoint(svc),
		Transport.DecodeSetUserRequest,
		Transport.EncodeResponse,
	)

	getRequestMetricsHandler := httptransport.NewServer(
		Transport.GetMetricsEndpoint(svc),
		Transport.DecodeMetricsRequest,
		Transport.EncodeResponse,
	)

	http.Handle("/getuser", getGetUserHandler)
	http.Handle("/setuser", getSetUserHandler)
	http.Handle("/createuser", getCreateUserHandler)
	http.Handle("/requestmetrics", getRequestMetricsHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":"+strconv.Itoa(serviceInfo.Port))
	logger.Log("err", http.ListenAndServe(":"+strconv.Itoa(serviceInfo.Port), nil))

}
