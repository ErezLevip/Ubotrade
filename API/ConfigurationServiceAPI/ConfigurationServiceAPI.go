package ConfigurationServiceAPI

import (
	"log"
	"reflect"
	"context"

	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/API/BaseServiceAPI"
	"github.com/erezlevip/Ubotrade/API/RegistryServiceAPI"
	"github.com/erezlevip/Ubotrade/Global"
)

const ConfigurationServiceKey = "Configuration"

type IConfigurationServiceAPI interface {
	GetConfiguration(ctx context.Context, req Global.ConfigurationRequest) (Global.ConfigurationResponse, error)
}

type ConfigurationServiceAPI struct {
	Base BaseServiceAPI.IBaseServiceAPI
}

func (svc ConfigurationServiceAPI) Make(ctx context.Context, basicConfigPath string) *ConfigurationServiceAPI {
	//get service information from registry service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	resp, err := registryService.GetService(ctx, Global.RegistryRequest{ServiceInformation: Global.ServiceInformation{ServiceName: ConfigurationServiceKey}})

	if err != nil {
		log.Println(err.Error())
	}
	//set the service information inside the response config service api instance
	configApi := ConfigurationServiceAPI{}
	configApi.Base = BaseServiceAPI.Make(ctx, basicConfigPath, ConfigurationServiceKey, &resp.ServiceInformation)
	return &configApi
}

func (svc *ConfigurationServiceAPI) GetConfiguration(ctx context.Context, req Global.ConfigurationRequest) (response Global.ConfigurationResponse, err error) {
	err = svc.Base.SendRequest(req, "getconfiguration", &response)
	return
}
