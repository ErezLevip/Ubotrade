package Transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"local/UbotTrade/Global"
	"local/UbotTrade/UserService/Service"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
)

func GetUserEndpoint(svc UserService.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.GetUserRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		data, err := svc.GetUser(ctx, req.UserId, req.DataType, req.ActiveOnly)
		if err != nil {
			return Global.LoginResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.GetUserResponse{Data: data, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func CreateUserEndpoint(svc UserService.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.CreateUserRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		data, err := svc.CreateUser(ctx, req.UserId, req.FirstName, req.LastName, req.Email)
		if err != nil {
			return Global.CreateUserResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.CreateUserResponse{Data: data, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func SetUserEndpoint(svc UserService.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.SetUserRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		err := svc.SetUser(ctx, req.UserId, req.DataType, req.Operation, req.Data)
		if err != nil {
			return Global.SetUserResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.SetUserResponse{Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetMetricsEndpoint(svc UserService.IUserService) endpoint.Endpoint {
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

func DecodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.CreateUserRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.GetUserRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeSetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.SetUserRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeTokenValidationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
