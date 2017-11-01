package Global

import (
	"crypto/md5"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type BaseServiceRequest struct {
}
type BaseServiceResponse struct {
	Status int `json:"status"`
}

type TokenRequest struct {
	ClientId string `json:"client_id"`
}
type GetUserRequest struct {
	UserId     string `json:"user_id"`
	DataType   string `json:"data_type"`
	ActiveOnly bool   `json:"active_only"`
}

type CreateUserRequest struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type SetUserRequest struct {
	UserId    string                 `json:"user_id"`
	DataType  string                 `json:"data_type"`
	Data      map[string]interface{} `json:"data"`
	Operation string                 `json:"operation"`
}

type CreateUserResponse struct {
	Base BaseServiceResponse    `json:"base"`
	Data map[string]interface{} `json:"data"`
}

type SetUserResponse struct {
	Base BaseServiceResponse `json:"base"`
}

type GetUserResponse struct {
	Base BaseServiceResponse      `json:"base"`
	Data []map[string]interface{} `json:"data"`
}

type LoginRequest struct {
	SessionId string `json:"session_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ClientId  string `json:"client_id"`
}

type LoginResponse struct {
	Base   BaseServiceResponse `json:"base"`
	UserId string              `json:"user_id"`
}

type TokenValidationRequest struct {
	Token string `json:"token"`
}
type TokenResponse struct {
	Base    BaseServiceResponse `json:"base"`
	Token   string              `json:"token"`
	IsValid bool                `json:"is_valid"`
}

type RegistryRequest struct {
	ServiceInformation ServiceInformation `json:"serviceinformation`
}

type RegistryServiceResponse struct {
	Base               BaseServiceResponse `json:"base"`
	ServiceInformation ServiceInformation  `json:"serviceinformation`
}

type BotProfitsRequest struct {
	Id        uuid.UUID `json:"id"`
	BotNumber int       `json:"bot_number"`
	Days      int       `json:"days"`
}

type GetAllActiveBotsRequest struct {
	UserId string `json:"user_id"`
}

type BotInformationRequest struct {
	Id             uuid.UUID `json:"id"`
	BotNumber      int       `json:"bot_number"`
	MaxResultCount int       `json:"max_result_count"`
}

type BotInformationResponse struct {
	Base BaseServiceResponse `json:"base"`
	Data []BotInformation    `json:"data"`
}

type GetAllBotsResponse struct {
	Base BaseServiceResponse `json:"base"`
	Data []GeneralBotInfo    `json:"data"`
}

type GeneralBotInfo struct {
	BotName   string    `json:"bot_name"`
	Id        uuid.UUID `json:"id"`
	BotNumber int       `json:"bot_number"`
}

type CreateBotRequest struct {
	TradingConfiguration TradingConfig `json:"trading_configuration"`
	UserId               string        `json:"user_id"`
}
type CreateBotResponse struct {
	Base BaseServiceResponse `json:"base"`
}

type BotTickerDataResponse struct {
	Base BaseServiceResponse  `json:"base"`
	Data []BotTickerDataModel `json:"data"`
}

type BotActivitiesResponse struct {
	Base BaseServiceResponse `json:"base"`
	Data []ActivityModel     `json:"data"`
}

type ConfigurationRequest struct {
	Key string `json:"key"`
}

type ConfigurationResponse struct {
	Base          BaseServiceResponse    `json:"base"`
	Configuration map[string]interface{} `json:"configuration"`
}

type MachineStateResponse struct {
	Data map[string]interface{} `json:"data"`
}

type ServiceMetricsResponse struct {
	Base BaseServiceResponse    `json:"base"`
	Data map[string]interface{} `json:"data"`
}

type ServiceConfiguration struct {
	Url string
}

type ServiceInformation struct {
	ServiceName     string    `json:"servicename"`
	Url             string    `json:"url"`
	Port            int       `json:"port`
	MachineName     string    `json:"machine_name"`
	LastHealthCheck time.Time `json:"last_health_check"`
}

func (serviceInfo ServiceInformation) GetHashCode(fixed bool) string {
	var comparableServiceInfo ServiceInformation
	if fixed {
		comparableServiceInfo = ServiceInformation{
			MachineName: serviceInfo.MachineName,
			Port:        serviceInfo.Port,
			ServiceName: serviceInfo.ServiceName,
			Url:         serviceInfo.Url,
		}
	} else {
		comparableServiceInfo = serviceInfo
	}
	md5 := md5.New()
	res, err := bson.Marshal(comparableServiceInfo)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	md5.Write(res)
	return string(md5.Sum(nil))
}

type Creds struct {
	username string
	password string
}

type TradingConfig struct {
	Stairs                         []Stair `json:"stairs"`
	Fallback                       float64 `json:"fallback"`
	BaseCommission                 float64 `json:"base_commission"`
	LastStairTimeout               string   `json:"last_stair_timeout"`
	TimeIterations                 string   `json:"time_iterations"`
	Currency                       string  `json:"currency"`
	PriceBlockForSlop              int     `json:"price_block_for_slop"`
	BuyOnPositiveSlop              bool    `json:"buy_on_positive_slop"`
	RetryTimeDuration              string   `json:"retry_time_duration"`
	PercentagePriceDifferenceOnBuy float64 `json:"percentage_price_difference_on_buy"`
	BotNumber                      int     `json:"bot_number"`
}

type Stair struct {
	Ratio float64 `json:"ratio"`
}

/*
type BotInfoRequest struct {
	Id uuid.UUID `json:"id"`
	BotNumber int `json:"bot_number"`
} */
