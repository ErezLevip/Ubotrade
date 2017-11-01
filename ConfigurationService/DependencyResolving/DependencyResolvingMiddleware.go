package DependencyResolving

import (
	"time"
	"local/UbotTrade/Global"
	"local/UbotTrade/DataHandlers/MongoDB"
	"log"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/Redis"
	"context"
	"local/UbotTrade/StaticConfiguration"
)

var LastConfigUpdate time.Time
var MongoDbConfiguration map[string]string
var RedisConfiguration DataHandlers.RedisConfiguration

type DependencyResolvingMiddleware struct {
	StaticConfigurationPath string
	Next ConfigurationService.IConfigurationService
}

func Make(next ConfigurationService.IConfigurationService, staticConfigPath string) (instance *DependencyResolvingMiddleware) {
	getDependenciesConfiguration(staticConfigPath)

	return &DependencyResolvingMiddleware{
		Next: next,
		StaticConfigurationPath:staticConfigPath,
	}
}

func GetRequestContext(ctx context.Context ,staticConfigPath string) (context.Context) {

	if (ctx == nil) {
		ctx = context.Background()
	}

	if (LastConfigUpdate.Add(30 * time.Minute).Before(time.Now())) {
		getDependenciesConfiguration(staticConfigPath)
	}

	//resolve handlers
	mongoHandler := &MongoDB.MongoHandler{}
	mongoHandler.Init(MongoDbConfiguration)
	ctx = context.WithValue(ctx, "MongoHandler", mongoHandler)
	redisHandler := DataHandlers.RedisHandler{}
	redisHandler.Init(RedisConfiguration)
	ctx = context.WithValue(ctx, "RedisHandler", redisHandler)
	//resolve configurations
	ctx = context.WithValue(ctx, "RedisConfiguration", RedisConfiguration)
	ctx = context.WithValue(ctx, "MongoConfiguration", MongoDbConfiguration)

	return ctx
}

func getDependenciesConfiguration(staticConfigPath string) {
	//load static configurations
	config, err := StaticConfiguration.ReadConfiguration(staticConfigPath) // make singleton

	MongoDbConfiguration = map[string]string{

		MongoDB.ConnectionStringKey: config[ConfigurationService.MongoConnectionStringKey].(string),
		MongoDB.DatabaseKey:         config[ConfigurationService.MongoDatabaseKey].(string),
	}

	RedisConfiguration = DataHandlers.RedisConfiguration{
		ConnectionString: config[ConfigurationService.RedisConnectionStringKey].(string),
		Credentials:      config[ConfigurationService.RedisCredentialsKey].(string),
		Db:               int(config[ConfigurationService.RedisDatabaseKey].(float64)),
	}

	if err != nil {
		log.Panicln(err.Error())
	}

	LastConfigUpdate = time.Now()
}
func (mw DependencyResolvingMiddleware) GetConfiguration(ctx context.Context, key string) (configuration map[string]interface{}, err error) {
	ctx = GetRequestContext(ctx,mw.StaticConfigurationPath)
	configuration, err = mw.Next.GetConfiguration(ctx, key)
	return
}

func (mw DependencyResolvingMiddleware) GetServiceMetrics() (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func (mw DependencyResolvingMiddleware) Init(staticConfigPath string, serviceInfo Global.ServiceInformation) *ConfigurationService.ConfigurationService {
	return mw.Next.Init(staticConfigPath, serviceInfo)
}
