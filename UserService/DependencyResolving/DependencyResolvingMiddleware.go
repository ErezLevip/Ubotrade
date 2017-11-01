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
	"local/UbotTrade/UserService/Service"
)

var LastConfigUpdate time.Time
var MongoDbConfiguration map[string]string
var RedisConfiguration DataHandlers.RedisConfiguration
var RabbitConfiguration RabbitMQ.RabbitConfiguration

type DependencyResolvingMiddleware struct {
	Next UserService.IUserService
}

func Make(next UserService.IUserService) (instance *DependencyResolvingMiddleware) {
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
		Topic: "User_Requests"}
	LastConfigUpdate = time.Now()
}

func (mw DependencyResolvingMiddleware) GetUser(ctx context.Context, userId string, dataType string, activeOnly bool) (data []map[string]interface{}, err error) {
	ctx = GetRequestContext(ctx)
	data, err = mw.Next.GetUser(ctx, userId, dataType, activeOnly)
	return
}

func (mw DependencyResolvingMiddleware) CreateUser(ctx context.Context, userId string, firstName string, lastName string, email string) (data map[string]interface{}, err error) {
	ctx = GetRequestContext(ctx)
	data, err = mw.Next.CreateUser(ctx, userId, firstName, lastName, email)
	return
}

func (mw DependencyResolvingMiddleware) SetUser(ctx context.Context, userId string, dataType string, operation string, data map[string]interface{}) (err error) {
	ctx = GetRequestContext(ctx)
	err = mw.Next.SetUser(ctx, userId, dataType, operation, data)
	return
}

func (mw DependencyResolvingMiddleware) Init(serviceInfo Global.ServiceInformation) *UserService.UserService {
	return mw.Next.Init(serviceInfo)
}

func (mw DependencyResolvingMiddleware) GetServiceMetrics() (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}
