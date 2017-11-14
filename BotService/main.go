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
	"github.com/erezlevip/Ubotrade/BotService/DependencyResolving"
	"github.com/erezlevip/Ubotrade/Global/ServiceHealth"
	"github.com/erezlevip/Ubotrade/BotService/Instrumenting"
	"github.com/erezlevip/Ubotrade/BotService/Logging"
	"github.com/erezlevip/Ubotrade/BotService/Service"
	"github.com/erezlevip/Ubotrade/BotService/Transport"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/Logger"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ubot",
		Subsystem: "bot_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "bot_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ubot",
		Subsystem: "bot_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	machineName, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
	}

	serviceInfo := Global.ServiceInformation{
		Port:            8087,
		ServiceName:     "Bot",
		Url:             "http://localhost",
		MachineName:     machineName,
		LastHealthCheck: time.Now(),
	}

	Logger.SetGlobalLogger()

	var svc BotService.IBotService
	svc = BotService.BotService{}.Init(serviceInfo)
	svc = DependencyResolving.Make(svc)
	svc = BotLogging.LoggingMiddleware{logger, svc}
	svc = BotInstrumenting.Make(requestCount, requestLatency, countResult, svc)

	ctx := DependencyResolving.GetRequestContext(context.Background())
	ServiceHealth.StartHealthTicker(ctx, serviceInfo)

	getBotInformationHandler := httptransport.NewServer(
		Transport.GetBotInformationEndpoint(svc),
		Transport.DecodeBotDataRequest,
		Transport.EncodeResponse,
	)

	getBotActivitiesHandler := httptransport.NewServer(
		Transport.GetLastActivitiesEndpoint(svc),
		Transport.DecodeBotDataRequest,
		Transport.EncodeResponse,
	)
	getAllActiveBots := httptransport.NewServer(
		Transport.GetAllActiveBotsEndpoint(svc),
		Transport.DecodeGetAllActiveBotsRequest,
		Transport.EncodeResponse,
	)

	getTickerDataHandler := httptransport.NewServer(
		Transport.GetBotTickerDataEndpoint(svc),
		Transport.DecodeBotDataRequest,
		Transport.EncodeResponse,
	)
	getBotProfitsHandler := httptransport.NewServer(
		Transport.GetBotProfitsDataEndpoint(svc),
		Transport.DecodeBotProfitsRequest,
		Transport.EncodeResponse,
	)

	createBotHandler := httptransport.NewServer(
		Transport.CreateBotEndpoint(svc),
		Transport.DecodeCreateBotRequest,
		Transport.EncodeResponse,
	)

	getRequestMetricsHandler := httptransport.NewServer(
		Transport.GetMetricsEndpoint(svc),
		Transport.DecodeMetricsRequest,
		Transport.EncodeResponse,
	)

	http.Handle("/getbotinformation", getBotInformationHandler)
	http.Handle("/getlastactivities", getBotActivitiesHandler)
	http.Handle("/getallactivebots", getAllActiveBots)
	http.Handle("/getbottickerdata", getTickerDataHandler)
	http.Handle("/getbotprofits", getBotProfitsHandler)
	http.Handle("/createnewbot", createBotHandler)
	http.Handle("/requestmetrics", getRequestMetricsHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":"+strconv.Itoa(serviceInfo.Port))
	logger.Log("err", http.ListenAndServe(":"+strconv.Itoa(serviceInfo.Port), nil))

}
