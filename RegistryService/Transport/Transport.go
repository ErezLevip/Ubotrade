package RegistryServiceTransport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"local/UbotTrade/Global"
	"local/UbotTrade/RegistryService/Service"
	"net/http"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
)

func RegisterEndpoint(svc RegistryService.IRegistryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.RegistryRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request, &req)

		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		err := svc.Register(ctx, req.ServiceInformation)
		if err != nil {
			return Global.RegistryServiceResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.RegistryServiceResponse{Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func DeRegisterEndpoint(svc RegistryService.IRegistryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.RegistryRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)

		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		err := svc.DeRegister(ctx, req.ServiceInformation)
		if err != nil {
			return Global.RegistryServiceResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.RegistryServiceResponse{Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetServiceEndpoint(svc RegistryService.IRegistryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.RegistryRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)

		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		serviceInformation, err := svc.GetService(ctx, req.ServiceInformation.ServiceName)
		if err != nil {
			return Global.RegistryServiceResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.RegistryServiceResponse{ServiceInformation: serviceInformation, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func DecodeRegistryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.RegistryRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
