package AuthenticationServiceAPI

import (
	"local/UbotTrade/API/BaseServiceAPI"
	"local/UbotTrade/API/RegistryServiceAPI"
	"local/UbotTrade/Global"
	"log"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"context"
)

const AuthenticationServiceKey = "Authentication"

type IAuthenticationServiceAPI interface {
	Login(req Global.LoginRequest) (response Global.LoginResponse, err error)
	GetToken(req Global.TokenRequest) (response Global.TokenResponse, err error)
	ValidateToken(req Global.TokenValidationRequest) (response Global.TokenResponse, err error)
}

type AuthenticationServiceAPI struct {
	Base *BaseServiceAPI.BaseServiceAPI
}

func (svc AuthenticationServiceAPI) Make(ctx context.Context, basicConfigPath string) *AuthenticationServiceAPI {
	//get service information from registry service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	resp, err := registryService.GetService(ctx, Global.RegistryRequest{ServiceInformation: Global.ServiceInformation{ServiceName: AuthenticationServiceKey}})

	if err != nil {
		log.Println(err.Error())
	}
	//set the service information inside the response config service api instance
	authAPI := AuthenticationServiceAPI{}
	authAPI.Base = BaseServiceAPI.Make(ctx, basicConfigPath, AuthenticationServiceKey, &resp.ServiceInformation)
	return &authAPI
}

// not in use at the moment
func (svc *AuthenticationServiceAPI) GetToken(req Global.TokenRequest) (response Global.TokenResponse, err error) {
	err = svc.Base.SendRequest(req, "gettoken", &response)
	return
}
// login and check the user session
func (svc *AuthenticationServiceAPI) Login(req Global.LoginRequest) (response Global.LoginResponse, err error) {
	err = svc.Base.SendRequest(req, "login", &response)
	return
}
// not in use at the moment
func (svc *AuthenticationServiceAPI) ValidateToken(req Global.TokenValidationRequest) (response Global.TokenResponse, err error) {
	err = svc.Base.SendRequest(req, "validatetoken", &response)
	return
}


