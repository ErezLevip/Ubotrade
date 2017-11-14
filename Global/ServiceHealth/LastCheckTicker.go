package ServiceHealth

import (
	"encoding/json"
	"log"
	"time"
	"context"

	"github.com/erezlevip/Ubotrade/DataHandlers/Redis"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/RegistryService/Service"
)

func StartHealthTicker(ctx context.Context, serviceInfo Global.ServiceInformation) {
	redisHandler := ctx.Value("RedisHandler").(DataHandlers.RedisHandler)

	serviceKey := RegistryService.ServicePrefix + "." + serviceInfo.ServiceName
	ticker := time.NewTicker(time.Second * 30)

	go func() {
		for t := range ticker.C {
			var services []Global.ServiceInformation
			servicesJson, err := redisHandler.Get(serviceKey)
			if err != nil {
				log.Println(err.Error())
			}

			if servicesJson != "" {
				err = json.Unmarshal([]byte(servicesJson), &services)
				if err != nil {
					log.Println(err.Error())
				}

				for i, service := range services {

					if serviceInfo.GetHashCode(true) == service.GetHashCode(true) {
						services[i].LastHealthCheck = t

						key := serviceInfo.GetHashCode(false)
						lock, err := redisHandler.Lock(key)
						if err != nil {
							log.Println(err.Error())
						}

						err = redisHandler.Set(serviceKey, services, 0)

						redisHandler.UnLock(lock)
						if err != nil {
							log.Println(time.Now(),"Failed to set health for key", serviceKey, err.Error())
						} else {
							log.Println(time.Now(),"Health status success", serviceKey, t)
						}
						break
					}
				}
			}
		}
	}()

}
