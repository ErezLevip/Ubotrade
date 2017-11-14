package RegistryService

import (
	"encoding/json"
	"errors"
	"github.com/erezlevip/Ubotrade/API/BaseServiceAPI"
	"github.com/erezlevip/Ubotrade/DataHandlers/Redis"
	"github.com/erezlevip/Ubotrade/Global"
	"log"
	"math/rand"
	"time"
	"context"
)

const ServicePrefix = "Services"
const MinPrecentOfFailuresPerService  = 50

const BasicConfigPath = "RegistryService/ServiceConfiguration.json"

type IRegistryService interface {
	Register(context.Context , Global.ServiceInformation) error
	DeRegister(context.Context, Global.ServiceInformation) error
	Init(config DataHandlers.RedisConfiguration) *RegistryService
	GetService(context.Context, string) (Global.ServiceInformation, error)
}

type RegistryService struct {
	redisHandler DataHandlers.RedisHandler
}

var RedisConfiguration DataHandlers.RedisConfiguration

func (svc RegistryService) Init(config DataHandlers.RedisConfiguration) *RegistryService {
	log.Println(time.Now(),"Starting Registry Service")
	RedisConfiguration = config
	return &RegistryService{
		redisHandler:DataHandlers.RedisHandler{},
	}
}

func (svc RegistryService) Register(ctx context.Context, serviceInfo Global.ServiceInformation) (err error) {
	svc.redisHandler.Init(RedisConfiguration)
	serviceKey := ServicePrefix + "." + serviceInfo.ServiceName

	var services []Global.ServiceInformation
	servicesJson, err := svc.redisHandler.Get(serviceKey)
	if err != nil {
		log.Println(err.Error())
	}
	if servicesJson == "" {
		services = make([]Global.ServiceInformation, 0)
		services = append(services, serviceInfo)
		err = svc.redisHandler.Set(serviceKey, services, 0)
	} else {
		err = json.Unmarshal([]byte(servicesJson), &services)
		if err != nil {
			log.Println(err.Error())
		}
		serviceExists := false
		for index := 0; index < len(services); index++ {
			if services[index].GetHashCode(true) == serviceInfo.GetHashCode(true) {
				serviceExists = true
				break
			}

		}
		if !serviceExists {
			services = append(services, serviceInfo)
			key := serviceInfo.GetHashCode(false)
			lock, err := svc.redisHandler.Lock(key)
			if err != nil {
				log.Println(err.Error())
			}
			err = svc.redisHandler.Set(serviceKey, services, 0)
			svc.redisHandler.UnLock(lock)
		}
	}

	return err
}

func (svc RegistryService) DeRegister(ctx context.Context, serviceInfo Global.ServiceInformation) (err error) {
	svc.redisHandler.Init(RedisConfiguration)
	serviceKey := ServicePrefix + "." + serviceInfo.ServiceName

	var services []Global.ServiceInformation

	servicesJson, err := svc.redisHandler.Get(serviceKey)
	if err != nil {
		log.Println(err.Error())
	}
	if servicesJson != "" {
		err = json.Unmarshal([]byte(servicesJson), &services)
		if err != nil {
			log.Println(err.Error())
		}
		serviceExists := false
		for index := 0; index <= len(services); index++ {
			if services[index].GetHashCode(true) == serviceInfo.GetHashCode(true) {
				services = append(services[:index], services[index+1:]...)
				serviceExists = true
				break
			}

		}
		if serviceExists {
			key := serviceInfo.GetHashCode(false)
			lock, err := svc.redisHandler.Lock(key)
			if err != nil {
				log.Println(err.Error())
			}
			err = svc.redisHandler.Set(serviceKey, services, 0)
			svc.redisHandler.UnLock(lock)
		}
	}
	return
}

func (svc RegistryService) GetService(ctx context.Context, serviceName string) (serviceResponse Global.ServiceInformation, err error) {
	svc.redisHandler.Init(RedisConfiguration)
	serviceKey := ServicePrefix + "." + serviceName

	var services []Global.ServiceInformation
	servicesJson, err := svc.redisHandler.Get(serviceKey)

	if servicesJson != "" {
		err = json.Unmarshal([]byte(servicesJson), &services)

		if services != nil && len(services) > 0 {
			instances := make([]Global.ServiceInformation, 0)
			for _, service := range services {
				if service.LastHealthCheck.Add(time.Second * time.Duration(30)).After(time.Now()) {
					baseInstance := BaseServiceAPI.Make(ctx,BasicConfigPath, "", &service)

					var metricsResponse Global.ServiceMetricsResponse
					metricsResponse, err = baseInstance.GetRequestMetrics()
					if metricsResponse.Data["SuccessRate"].(float64) > MinPrecentOfFailuresPerService {
						log.Println(time.Now(),serviceName, "successrate is ", metricsResponse.Data["SuccessRate"])
						instances = append(instances, service)
					}
				}
			}
			if len(instances) > 0 {
				index := rand.Intn(len(instances))
				serviceResponse = instances[index]
				return
			} else {
				log.Println(time.Now(),serviceName," has 0 valid instances")
				err = errors.New("Failed to get an instance with a high success rate")
			}
		}
	} else {
		err = errors.New("Failed to get an instance with a high success rate")
	}
	return
}
/*
func getServiceMetrics(serviceInfo Global.ServiceInformation) (metricsResponse *Global.ServiceMetricsResponse, err error) {
	api := BaseServiceAPI.BaseServiceAPI{}
	err = api.SendRequestWithUrl(nil, "requestmetrics", metricsResponse, serviceInfo.Url+":"+strconv.Itoa(serviceInfo.Port))
	return
}

var ErrEmpty = errors.New("empty string")
*/