package Handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/erezlevip/Ubotrade/API/BotServiceAPI"
	"github.com/erezlevip/Ubotrade/Global"
)

type BotRequestModel struct {
	BotNumber int `json:"botNumber"`
}

type BotProfitsRequestModel struct {
	BotNumber int `json:"botNumber"`
	Days      int `json:"days"`
}

type DashBoardBotInfo struct {
}

type BotInfoHandler interface {
	GetBotActivity(w http.ResponseWriter, req *http.Request, ctx context.Context)
	GetBotInformation(w http.ResponseWriter, req *http.Request, ctx context.Context)
	GetLastActivities(w http.ResponseWriter, req *http.Request, ctx context.Context)
	GetBotProfits(w http.ResponseWriter, req *http.Request, ctx context.Context)
	CreateNewBot(w http.ResponseWriter, req *http.Request, ctx context.Context)
	CancelBot(w http.ResponseWriter, req *http.Request, ctx context.Context)
	GetAllActiveBots(w http.ResponseWriter, req *http.Request, ctx context.Context)
	GetBotTickerData(w http.ResponseWriter, req *http.Request, ctx context.Context)
}

func BotInfoHandlerMake() BotInfoHandler {
	return &DashBoardBotInfo{}
}

func (infoHandler *DashBoardBotInfo) GetBotActivity(w http.ResponseWriter, req *http.Request, ctx context.Context) {

}

func (infoHandler *DashBoardBotInfo) GetBotInformation(w http.ResponseWriter, req *http.Request, ctx context.Context) {

	var requestModel BotRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)
	log.Println(requestModel)
	if (requestModel != BotRequestModel{}) {
		botSvc := ctx.Value(reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)

		botInfoResponse, err := botSvc.GetBotInformation(Global.BotInformationRequest{
			BotNumber: requestModel.BotNumber,
		})

		log.Println(botInfoResponse)

		if err != nil {
			log.Println(err.Error())
		}
		if len(botInfoResponse.Data) > 0 {
			bot := botInfoResponse.Data[0]
			botNumber := bot.Configuration.BotNumber
			resp := make(map[string]interface{})
			resp["Currency"] = bot.Configuration.Currency
			resp["BotNumber"] = botNumber
			resp["Amount"] = bot.Amount
			resp["BotName"] = bot.Name
			json.NewEncoder(w).Encode(resp)
		} else {
			log.Println(time.Now(), "botInfoResponse is empty")
		}
	}
}

func (infoHandler *DashBoardBotInfo) GetLastActivities(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	botSvc := ctx.Value(reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)

	var requestModel BotRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)
	if (requestModel != BotRequestModel{}) {
		res, err := botSvc.GetLastActivities(Global.BotInformationRequest{
			BotNumber:      requestModel.BotNumber,
			MaxResultCount: 10,
		})
		if err != nil {
			log.Println(err.Error())
		}

		resp := make([]map[string]interface{}, 0)
		index := 1
		for _, activity := range res.Data {
			if (activity.Id != uuid.UUID{}) {
				activityData := make(map[string]interface{})
				activityData["Id"] = activity.Id
				activityData["ActivityPrice"] = activity.ActivityPrice
				activityData["ActivityType"] = activity.ActivityType
				activityData["ActualAmountCurrency"] = activity.ActualAmountCurrency
				activityData["ActualAmountUSD"] = activity.ActualAmountUSD
				activityData["BotId"] = activity.BotId
				activityData["PriceDifference"] = activity.PriceDifference
				activityData["Index"] = index
				activityData["TimeStamp"] = activity.TimeStamp
				resp = append(resp, activityData)
				index++
			}
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func (infoHandler *DashBoardBotInfo) GetBotProfits(w http.ResponseWriter, req *http.Request, ctx context.Context) {

	botSvc := ctx.Value(reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)

	var requestModel BotProfitsRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)
	if (requestModel != BotProfitsRequestModel{}) {
		res, err := botSvc.GetBotProfits(Global.BotProfitsRequest{
			BotNumber: requestModel.BotNumber,
			Days:      requestModel.Days,
		})
		if err != nil {
			log.Println(err.Error())
		}
		resp := make([]float64, 0)
		for _, activity := range res.Data {
			if (activity.Id != uuid.UUID{}) {
				resp = append(resp, activity.ActualAmountUSD)
			}
		}
		json.NewEncoder(w).Encode(resp)
	}
}

// not implemented on the clientside yet, server logic already implemented on bot service
func (infoHandler *DashBoardBotInfo) CreateNewBot(w http.ResponseWriter, req *http.Request, ctx context.Context) {

}

// not implemented on the clientside yet, server logic already implemented on bot service
func (infoHandler *DashBoardBotInfo) CancelBot(w http.ResponseWriter, req *http.Request, ctx context.Context) {

}

// get all the active bots of the logged in account
func (infoHandler *DashBoardBotInfo) GetAllActiveBots(w http.ResponseWriter, req *http.Request, ctx context.Context) {

	botSvc := ctx.Value(reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)

	//the app will use a testing id so other people will be able to see the bots
	userId := "118103040085940455572"

	allBots, err := botSvc.GetAllActiveBots(Global.GetAllActiveBotsRequest{
		//UserId: ctx.Value("UserId").(string),
		UserId: userId,
	})
	if err != nil {
		log.Println(err.Error())
	}

	resp := make([]map[string]interface{}, 0)
	for _, bot := range allBots.Data {
		botData := make(map[string]interface{})
		botData["Id"] = bot.Id
		botData["BotNumber"] = bot.BotNumber
		botData["BotName"] = bot.BotName
		resp = append(resp, botData)
	}
	json.NewEncoder(w).Encode(resp)

}

// get the real time ticker data of the current bot
func (infoHandler *DashBoardBotInfo) GetBotTickerData(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	botSvc := ctx.Value(reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)

	var requestModel BotRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)
	if (requestModel != BotRequestModel{}) {
		res, err := botSvc.GetBotTickerData(Global.BotInformationRequest{
			BotNumber:      requestModel.BotNumber,
			MaxResultCount: 10,
		})
		if err != nil {
			log.Println(err.Error())
		}
		resp := make([]map[string]interface{}, 0)
		for _, ticker := range res.Data {
			tickerData := make(map[string]interface{})
			tickerData["Currency"] = ticker.Currency
			tickerData["Action"] = ticker.Action
			tickerData["H"] = ticker.H
			tickerData["P"] = ticker.P
			tickerData["Stairs"] = ticker.Stairs
			resp = append(resp, tickerData)
		}
		json.NewEncoder(w).Encode(resp)
	}

}
