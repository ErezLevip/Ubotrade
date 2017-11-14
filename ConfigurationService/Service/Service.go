package ConfigurationService

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"context"
	"time"
	"encoding/json"

	"github.com/erezlevip/Ubotrade/DataHandlers/Redis"
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/API/RegistryServiceAPI"
	"github.com/erezlevip/Ubotrade/DataHandlers/MongoDB"
	"github.com/erezlevip/Ubotrade/Global"
)

const ServicePrefix = "Services"
const ConfigurationsCollection = "Configurations"
const ConnectionStringsConfigurationKey = "ConnectionStrings"

const RabbitMQConnectionStringKey = "rabbitmqconnectionstring"

const RedisConnectionStringKey = "redisconnectionstring"
const RedisDatabaseKey = "redisdb"
const RedisCredentialsKey = "rediscredentials"

const MongoConnectionStringKey = "mongodbconnectionstring"
const MongoDatabaseKey = "mongodatabase"
const RedisConfigurationsPrefix = "Configurations"

type IConfigurationService interface {
	GetConfiguration(ctx context.Context, key string) (map[string]interface{}, error)
	GetServiceMetrics() (metrics map[string]interface{}, err error)
	Init(staticConfigPath string, serviceInfo Global.ServiceInformation) *ConfigurationService
}

type ConfigurationService struct {
}

func (svc ConfigurationService) Init(staticConfigPath string, serviceInfo Global.ServiceInformation) *ConfigurationService {
	log.Println(time.Now(), "Starting Configuration Service")

	configSvc := ConfigurationService{}

	ctx := context.Background()

	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	_, err := registryService.Register(ctx, Global.RegistryRequest{ServiceInformation: serviceInfo})
	if (err != nil) {
		log.Fatal(err.Error())
	}
	return &configSvc
}

func (svc ConfigurationService) GetConfiguration(ctx context.Context, key string) (configuration map[string]interface{}, err error) {

	redisHandler := ctx.Value("RedisHandler").(DataHandlers.RedisHandler)

	var cachedConfig string
	configKey := RedisConfigurationsPrefix + "." + key
	cachedConfig, err = redisHandler.Get(key)
	if (cachedConfig == "") {
		mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
		configuration, err = mongoHandler.FindFirst(ConfigurationsCollection, bson.M{"Key": key})
		var jdata []byte
		jdata, err = json.Marshal(configuration)
		redisHandler.Set(configKey, string(jdata), time.Duration(30)*time.Minute)
	} else {
		err = json.Unmarshal([]byte(cachedConfig), &configuration)
	}

	return
}

func (svc ConfigurationService) GetServiceMetrics() (outputMetrics map[string]interface{}, err error) {
	return make(map[string]interface{}), nil
}
