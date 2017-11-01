package RegistryServiceAPI

import (
	"local/UbotTrade/API/BaseServiceAPI"
	"local/UbotTrade/Global"
	"log"
	"context"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
	"strings"
	"strconv"
)

const RegistryConfigKey = "RegistryService"

type IRegistryServiceAPI interface {
	Register(context.Context, Global.RegistryRequest) (Global.RegistryServiceResponse, error)
	DeRegister(context.Context, Global.RegistryRequest) (Global.RegistryServiceResponse, error)
	GetService(context.Context, Global.RegistryRequest) (Global.RegistryServiceResponse, error)
	GetMetrics(ctx context.Context) (interface{}, error)
}

type RegistryServiceAPI struct {
	Base BaseServiceAPI.BaseServiceAPI
}

func (svc RegistryServiceAPI) Make(ctx context.Context, basicConfigPath string) *RegistryServiceAPI {
	newService := RegistryServiceAPI{}

	//load the static config file
	BaseServiceAPI.LoadConfig(basicConfigPath)

	// registry service's basic information since it doesn't perform registration
	registryUrl := BaseServiceAPI.BasicConfig[RegistryConfigKey].(string)
	var registryUrlSplit = strings.Split(registryUrl, ":")
	port, err := strconv.Atoi(registryUrlSplit[2])
	if (err != nil) {
		log.Panic(err.Error())
	}
	machineName := strings.Split(registryUrlSplit[1], "//")

	//this data is passed through a header called hops
	newService.Base = *BaseServiceAPI.Make(ctx, basicConfigPath, RegistryConfigKey, &Global.ServiceInformation{
		ServiceName: RegistryConfigKey,
		Url:         registryUrl,
		Port:        port,
		MachineName: machineName[1],
	})
	return &newService
}

func (svc *RegistryServiceAPI) Register(ctx context.Context, req Global.RegistryRequest) (response Global.RegistryServiceResponse, err error) {
	requestMetadata := BaseServiceAPI.HandleRequestMetadata(ctx)
	GeneralMicroserviceBehavior.RegisterMicroserviceHop(requestMetadata, *svc.Base.ServiceInformation)
	err = svc.Base.SendRequestWithUrl(requestMetadata, req, "register", &response, BaseServiceAPI.BasicConfig[RegistryConfigKey].(string))
	return
}

func (svc *RegistryServiceAPI) DeRegister(ctx context.Context, req Global.RegistryRequest) (response Global.RegistryServiceResponse, err error) {
	requestMetadata := BaseServiceAPI.HandleRequestMetadata(ctx)
	GeneralMicroserviceBehavior.RegisterMicroserviceHop(requestMetadata, *svc.Base.ServiceInformation)
	err = svc.Base.SendRequestWithUrl(requestMetadata, req, "deregister", &response, BaseServiceAPI.BasicConfig[RegistryConfigKey].(string))
	return
}
func (svc *RegistryServiceAPI) GetService(ctx context.Context, req Global.RegistryRequest) (response Global.RegistryServiceResponse, err error) {
	requestMetadata := BaseServiceAPI.HandleRequestMetadata(ctx)
	GeneralMicroserviceBehavior.RegisterMicroserviceHop(requestMetadata, *svc.Base.ServiceInformation)
	err = svc.Base.SendRequestWithUrl(requestMetadata, req, "getservice", &response, BaseServiceAPI.BasicConfig[RegistryConfigKey].(string))
	return
}
