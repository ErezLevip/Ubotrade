package BotServiceAPI

import (
	"log"
	"reflect"
	"context"
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/API/BaseServiceAPI"
	"github.com/erezlevip/Ubotrade/API/RegistryServiceAPI"
	"github.com/erezlevip/Ubotrade/Global"
)

const BotServiceKey = "Bot"

type IBotServiceAPI interface {
	GetBotInformation(Global.BotInformationRequest) (Global.BotInformationResponse, error)
	GetAllActiveBots() (Global.GetAllBotsResponse, error)
	GetLastActivities(Global.BotInformationRequest) (Global.BotActivitiesResponse, error)
	GetBotTickerData(Global.BotInformationRequest) (Global.BotTickerDataResponse, error)
	GetBotProfits(req Global.BotProfitsRequest) (response Global.BotActivitiesResponse, err error)
	CreateNewBot(Global.CreateBotRequest) (Global.CreateBotResponse, error)
}

type BotServiceAPI struct {
	Base *BaseServiceAPI.BaseServiceAPI
}

func (svc BotServiceAPI) Make(ctx context.Context, basicConfigPath string) *BotServiceAPI {
	//get service information from registry service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	resp, err := registryService.GetService(ctx, Global.RegistryRequest{ServiceInformation: Global.ServiceInformation{ServiceName: BotServiceKey}})

	if err != nil {
		log.Println(err.Error())
	}
	//set the service information inside the response config service api instance
	botApi := BotServiceAPI{}
	botApi.Base = BaseServiceAPI.Make(ctx, basicConfigPath, BotServiceKey, &resp.ServiceInformation)
	return &botApi
}

func (svc *BotServiceAPI) GetBotInformation(req Global.BotInformationRequest) (response Global.BotInformationResponse, err error) {
	err = svc.Base.SendRequest(req, "getbotinformation", &response)
	return
}

func (svc *BotServiceAPI) GetLastActivities(req Global.BotInformationRequest) (response Global.BotActivitiesResponse, err error) {
	err = svc.Base.SendRequest(req, "getlastactivities", &response)
	return
}

func (svc *BotServiceAPI) GetBotTickerData(req Global.BotInformationRequest) (response Global.BotTickerDataResponse, err error) {
	err = svc.Base.SendRequest(req, "getbottickerdata", &response)
	return
}

func (svc *BotServiceAPI) GetBotProfits(req Global.BotProfitsRequest) (response Global.BotActivitiesResponse, err error) {
	err = svc.Base.SendRequest(req, "getbotprofits", &response)
	return
}

func (svc *BotServiceAPI) CreateNewBot(req Global.CreateBotRequest) (response Global.CreateBotResponse, err error) {
	err = svc.Base.SendRequest(req, "createnewbot", &response)
	return
}
func (svc *BotServiceAPI) GetAllActiveBots(req Global.GetAllActiveBotsRequest) (response Global.GetAllBotsResponse, err error) {
	err = svc.Base.SendRequest(req, "getallactivebots", &response)
	return
}

