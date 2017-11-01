package DependencyResolving

import (
	"time"
	"local/UbotTrade/Global"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"local/UbotTrade/API/ConfigurationServiceAPI"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"log"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/Redis"
	"context"
	"local/UbotTrade/BotService/Service"
	"github.com/nu7hatch/gouuid"
)

var LastConfigUpdate time.Time
var MongoDbConfiguration map[string]string
var RedisConfiguration DataHandlers.RedisConfiguration
var RabbitConfiguration RabbitMQ.RabbitConfiguration

type DependencyResolvingMiddleware struct {
	Next BotService.IBotService
}

func Make(next BotService.IBotService) (instance *DependencyResolvingMiddleware) {
	getDependenciesConfiguration(context.Background())

	return &DependencyResolvingMiddleware{
		Next: next,
	}
}

func GetRequestContext(ctx context.Context) (context.Context) {

	if (ctx == nil) {
		ctx = context.Background()
	}

	if (LastConfigUpdate.Add(30 * time.Minute).Before(time.Now())) {
		getDependenciesConfiguration(ctx)
	}

	//resolve handlers
	mongoHandler := &MongoDB.MongoHandler{}
	mongoHandler.Init(MongoDbConfiguration)
	ctx = context.WithValue(ctx, "MongoHandler", mongoHandler)
	redisHandler := DataHandlers.RedisHandler{}
	redisHandler.Init(RedisConfiguration)
	ctx = context.WithValue(ctx, "RedisHandler", redisHandler)
	rabbitHandler := &RabbitMQ.RabbitHandler{}
	rabbitHandler.Init(RabbitConfiguration)
	ctx = context.WithValue(ctx, "RabbitHandler", rabbitHandler)
	//resolve configurations
	ctx = context.WithValue(ctx, "RedisConfiguration", RedisConfiguration)
	ctx = context.WithValue(ctx, "MongoConfiguration", MongoDbConfiguration)
	ctx = context.WithValue(ctx, "RabbitConfiguration", RabbitConfiguration)

	return ctx
}

func getDependenciesConfiguration(ctx context.Context) {
	//load static configurations
	configService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(ConfigurationServiceAPI.ConfigurationServiceAPI{})).(*ConfigurationServiceAPI.ConfigurationServiceAPI)
	connectionsResponse, err := configService.GetConfiguration(ctx, Global.ConfigurationRequest{Key: ConfigurationService.ConnectionStringsConfigurationKey})

	if err != nil {
		log.Println(err.Error())
	}
	//create rabbitMQ and Mongo Db handlers with their configurations
	connectionStringsConfig := connectionsResponse.Configuration["Data"].(map[string]interface{})

	MongoDbConfiguration = map[string]string{
		MongoDB.ConnectionStringKey: connectionStringsConfig[ConfigurationService.MongoConnectionStringKey].(string),
		MongoDB.DatabaseKey:         connectionStringsConfig[ConfigurationService.MongoDatabaseKey].(string),
	}

	RedisConfiguration = DataHandlers.RedisConfiguration{
		ConnectionString: connectionStringsConfig[ConfigurationService.RedisConnectionStringKey].(string),
		Credentials:      connectionStringsConfig[ConfigurationService.RedisCredentialsKey].(string),
		Db:               int(connectionStringsConfig[ConfigurationService.RedisDatabaseKey].(float64)),
	}

	RabbitConfiguration = RabbitMQ.RabbitConfiguration{ConnectionString: connectionStringsConfig[ConfigurationService.RabbitMQConnectionStringKey].(string),
		Topic: "Bot_Requests"}
	LastConfigUpdate = time.Now()
}

func (mw DependencyResolvingMiddleware) GetBotInformation(ctx context.Context, Id uuid.UUID, BotNumber int) (bots Global.BotInformation, err error) {
	ctx = GetRequestContext(ctx)
	bots, err = mw.Next.GetBotInformation(ctx, Id, BotNumber)
	return
}
func (mw DependencyResolvingMiddleware) GetLastActivities(ctx context.Context, Id uuid.UUID, BotNumber int, MaxResultCount int) (activities []Global.ActivityModel, err error) {
	ctx = GetRequestContext(ctx)
	activities, err = mw.Next.GetLastActivities(ctx, Id, BotNumber, MaxResultCount)
	return
}

func (mw DependencyResolvingMiddleware) GetAllActiveBots(ctx context.Context, userId string) (bots []Global.GeneralBotInfo, err error) {
	ctx = GetRequestContext(ctx)
	bots, err = mw.Next.GetAllActiveBots(ctx, userId)
	return
}

func (mw DependencyResolvingMiddleware) GetBotTickerData(ctx context.Context,  id uuid.UUID, botNumber int, maxResultCount int) (prices []Global.BotTickerDataModel, err error) {
	ctx = GetRequestContext(ctx)
	prices, err = mw.Next.GetBotTickerData(ctx , id, botNumber, maxResultCount)
	return
}
func (mw DependencyResolvingMiddleware) GetBotProfits(ctx context.Context,  id uuid.UUID, botNumber int, days int) (prices []Global.ActivityModel, err error) {
	ctx = GetRequestContext(ctx)
	prices, err = mw.Next.GetBotProfits(ctx, id, botNumber, days)
	return
}

func (mw DependencyResolvingMiddleware) CreateNewBot(ctx context.Context, tradingConfiguration Global.TradingConfig, UserId string) (err error) {
	ctx = GetRequestContext(ctx)
	err = mw.Next.CreateNewBot(ctx, tradingConfiguration, UserId)
	return
}

func (mw DependencyResolvingMiddleware) Init(serviceInfo Global.ServiceInformation) *BotService.BotService {
	return mw.Next.Init(serviceInfo)
}

func (mw DependencyResolvingMiddleware) GetServiceMetrics() (metrics map[string]interface{}, err error) {
	return
}