package Transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/erezlevip/Ubotrade/BotService/Service"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/Global/GeneralMicroserviceBehavior"
)

func GetBotInformationEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.BotInformationRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		bot, err := svc.GetBotInformation(ctx, req.Id, req.BotNumber)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.BotInformationResponse{Data: []Global.BotInformation{bot}, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}
func GetLastActivitiesEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.BotInformationRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		activities, err := svc.GetLastActivities(ctx, req.Id, req.BotNumber, req.MaxResultCount)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.BotActivitiesResponse{Data: activities, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetAllActiveBotsEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.GetAllActiveBotsRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		bots, err := svc.GetAllActiveBots(ctx, req.UserId)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.GetAllBotsResponse{Data: bots, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}
func GetBotTickerDataEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.BotInformationRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		tickerData, err := svc.GetBotTickerData(ctx, req.Id, req.BotNumber, req.MaxResultCount)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.BotTickerDataResponse{Data: tickerData, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}
func GetBotProfitsDataEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.BotProfitsRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		profits, err := svc.GetBotProfits(ctx, req.Id, req.BotNumber, req.Days)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}

		return Global.BotActivitiesResponse{Data: profits, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func CreateBotEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.CreateBotRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		err := svc.CreateNewBot(ctx, req.TradingConfiguration,req.UserId)
		if err != nil {
			return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.BotInformationResponse{Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetMetricsEndpoint(svc BotService.IBotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		metrics, err := svc.GetServiceMetrics()
		if err != nil {
			return Global.ServiceMetricsResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.ServiceMetricsResponse{Data: metrics, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func DecodeMetricsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request map[string]interface{}
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeGetAllActiveBotsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.GetAllActiveBotsRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeBotDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.BotInformationRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeBotProfitsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.BotProfitsRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeCreateBotRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.CreateBotRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

/*
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}*/
