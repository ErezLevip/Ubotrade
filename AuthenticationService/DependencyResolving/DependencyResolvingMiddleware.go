package DependencyResolving

import (
	"time"
	"local/UbotTrade/Global"
	"local/UbotTrade/AuthenticationService/Service"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"local/UbotTrade/API/ConfigurationServiceAPI"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"log"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/Redis"
	"context"
)

var LastConfigUpdate time.Time
var MongoDbConfiguration map[string]string
var RedisConfiguration DataHandlers.RedisConfiguration
var RabbitConfiguration RabbitMQ.RabbitConfiguration

type DependencyResolvingMiddleware struct {
	Next AuthenticationService.IAuthenticationService
}

// this middleware will resolve the service's critical dependencies once every request.

func Make(next AuthenticationService.IAuthenticationService) (instance *DependencyResolvingMiddleware) {
	getDependenciesConfiguration(context.Background())

	return &DependencyResolvingMiddleware{
		Next: next,
	}
}
// creating the context out of the Background context or if passed inserting all the data handlers to it
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
// getting all the configurations from configuration service
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
		Topic: "Authentication_Requests"}
	LastConfigUpdate = time.Now()
}

func (mw DependencyResolvingMiddleware) Login(ctx context.Context, clientId string, firstName string, lastName string, email string, sessionId string) (userId string, err error) {
	ctx = GetRequestContext(ctx)
	userId, err = mw.Next.Login(ctx, clientId, firstName, lastName, email, sessionId)
	return
}
func (mw DependencyResolvingMiddleware) GetToken(ctx context.Context, clientId string) (token string, err error) {
	ctx = GetRequestContext(ctx)
	token, err = mw.Next.GetToken(ctx, clientId)
	return
}
func (mw DependencyResolvingMiddleware) ValidateToken(ctx context.Context, token string) (IsValid bool, err error) {
	ctx = GetRequestContext(ctx)
	IsValid, err = mw.Next.ValidateToken(ctx, token)
	return
}

func (mw DependencyResolvingMiddleware) Init(serviceInfo Global.ServiceInformation) *AuthenticationService.AuthenticationService {
	return mw.Next.Init(serviceInfo)
}

func (mw DependencyResolvingMiddleware) GetServiceMetrics() (metrics map[string]interface{}, err error) {
	return
}
