package DependencyResolving

import (
	"time"
	"local/UbotTrade/Global"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/DataHandlers/Redis"
	"context"
	"local/UbotTrade/RegistryService/Service"
)

var LastConfigUpdate time.Time
var MongoDbConfiguration map[string]string
var RedisConfiguration DataHandlers.RedisConfiguration
var RabbitConfiguration RabbitMQ.RabbitConfiguration

type DependencyResolvingMiddleware struct {
	Next RegistryService.IRegistryService
}

func Make(next RegistryService.IRegistryService) (instance *DependencyResolvingMiddleware) {
	getDependenciesConfiguration(context.Background())

	return &DependencyResolvingMiddleware{
		Next: next,
	}
}

func GetRequestContext(ctx context.Context) (context.Context) {

	if (ctx == nil) {
		ctx = context.Background()
	}

	/*if (LastConfigUpdate.Add(30 * time.Minute).Before(time.Now())) {
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
*/
	return ctx
}

func getDependenciesConfiguration(ctx context.Context) {
	/*//load static configurations
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
	LastConfigUpdate = time.Now()*/
}

func (mw DependencyResolvingMiddleware) Register(ctx context.Context, serviceInfo Global.ServiceInformation) (err error) {
	ctx = GetRequestContext(ctx)
	 err = mw.Next.Register(ctx, serviceInfo)
	return
}

func (mw DependencyResolvingMiddleware) DeRegister(ctx context.Context, serviceInfo Global.ServiceInformation) (err error) {
	ctx = GetRequestContext(ctx)
	err = mw.Next.DeRegister(ctx, serviceInfo)
	return
}

func (mw DependencyResolvingMiddleware) GetService(ctx context.Context, serviceName string) (serviceResponse Global.ServiceInformation, err error) {
	ctx = GetRequestContext(ctx)
	serviceResponse, err = mw.Next.GetService(ctx, serviceName)
	return
}
/*
func getServiceMetrics(serviceInfo Global.ServiceInformation) (metricsResponse *Global.ServiceMetricsResponse, err error) {
	api := BaseServiceAPI.BaseServiceAPI{}
	err = api.SendRequestWithUrl(nil, "requestmetrics", metricsResponse, serviceInfo.Url+":"+strconv.Itoa(serviceInfo.Port))
	return
}*/

func (mw DependencyResolvingMiddleware) Init(config DataHandlers.RedisConfiguration) *RegistryService.RegistryService {
	return mw.Next.Init(config)
}

func (mw DependencyResolvingMiddleware) GetServiceMetrics() (metrics map[string]interface{}, err error) {
	return
}