package Transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"local/UbotTrade/Global"
	"local/UbotTrade/AuthenticationService/Service"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
)

func LoginEndpoint(svc AuthenticationService.IAuthenticationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.LoginRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		userId, err := svc.Login(ctx, req.ClientId, req.FirstName, req.LastName, req.Email, req.SessionId)
		if err != nil {
			return Global.LoginResponse{Base: Global.BaseServiceResponse{Status: http.StatusUnauthorized}}, err
		}
		return Global.LoginResponse{UserId: userId, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}
func GetTokenEndpoint(svc AuthenticationService.IAuthenticationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.TokenRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		token, err := svc.GetToken(ctx, req.ClientId)
		if err != nil {
			return Global.TokenResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.TokenResponse{Token: token, IsValid: true, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}
func ValidateTokenEndpoint(svc AuthenticationService.IAuthenticationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req Global.TokenValidationRequest
		GeneralMicroserviceBehavior.GetRequestFromGenericRequest(request,&req)
		ctx = GeneralMicroserviceBehavior.GetContextFromGenericRequest(ctx, request)
		isValid, err := svc.ValidateToken(ctx, req.Token)
		if err != nil {
			return Global.TokenResponse{Base: Global.BaseServiceResponse{Status: http.StatusInternalServerError}}, err
		}
		return Global.TokenResponse{IsValid: isValid, Token: req.Token, Base: Global.BaseServiceResponse{Status: http.StatusOK}}, nil
	}
}

func GetMetricsEndpoint(svc AuthenticationService.IAuthenticationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
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

func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.LoginRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.TokenRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func DecodeTokenValidationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Global.TokenValidationRequest
	return GeneralMicroserviceBehavior.DecodeRequestWithMetadata(request, r)
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
