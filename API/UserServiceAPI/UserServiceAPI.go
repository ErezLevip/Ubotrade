package UserServiceAPI

import (
	"local/UbotTrade/API/BaseServiceAPI"
	"local/UbotTrade/API/RegistryServiceAPI"
	"local/UbotTrade/Global"
	"log"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"context"
)

const UserServiceKey = "User"

type IUserServiceAPI interface {
	GetUser(req Global.GetUserRequest) (response Global.GetUserResponse, err error)
	SetUser(req Global.GetUserRequest) (response Global.GetUserResponse, err error)
	CreateUser(req Global.GetUserRequest) (response Global.GetUserResponse, err error)
}

type UserServiceAPI struct {
	Base *BaseServiceAPI.BaseServiceAPI
}

func (svc UserServiceAPI) Make(ctx context.Context, basicConfigPath string) *UserServiceAPI {
	//get service information from registry service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	resp, err := registryService.GetService(ctx, Global.RegistryRequest{ServiceInformation: Global.ServiceInformation{ServiceName: UserServiceKey}})

	if err != nil {
		log.Println(err.Error())
	}
	//set the service information inside the response config service api instance
	authAPI := UserServiceAPI{}
	authAPI.Base = BaseServiceAPI.Make(ctx, basicConfigPath, UserServiceKey, &resp.ServiceInformation)
	return &authAPI
}

func (svc *UserServiceAPI) GetUser(req Global.GetUserRequest) (response Global.GetUserResponse, err error) {
	err = svc.Base.SendRequest(req, "getuser", &response)
	return
}

func (svc *UserServiceAPI) SetUser(req Global.SetUserRequest) (response Global.SetUserResponse, err error) {
	err = svc.Base.SendRequest(req, "setuser", &response)
	return
}

func (svc *UserServiceAPI) CreateUser(req Global.CreateUserRequest) (response Global.CreateUserResponse, err error) {
	err = svc.Base.SendRequest(req, "createuser", &response)
	return
}
