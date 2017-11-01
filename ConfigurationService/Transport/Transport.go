package ConfigurationServiceTransport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/Global"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
)

func GetConfigurationEndpoint(svc ConfigurationService.IConfigurationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.ConfigurationRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request, &req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		config, err := svc.GetConfiguration(ctx, req.Key)
		if err != nil {
			return Global.ConfigurationResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.ConfigurationResponse{Configuration: config, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetMetricsEndpoint(svc ConfigurationService.IConfigurationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request).(Global.ConfigurationRequest)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		metrics, err := svc.GetServiceMetrics()
		if err != nil {
			return Global.ServiceMetricsResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.ServiceMetricsResponse{Data: metrics, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func DecodeConfigurationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.ConfigurationRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeMetricsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
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
