# Ubotrade

Ubotrade is a microservices based crypto currency trader.

I created the project to experience the language and to have a deeper understanding of some patterns such as middlewares and SOA.

# The project is based on 2 core microservices:

Registry service
Configuration service

Each one of these services is a critical pillar in the system, the registry service is in charge of service discovery while the configuration feed the relevant configuration for each service/module.

# The rest of the services are:

Authentication service -> authenticate/authorize each request to the application.
Bot service -> get update delete bots information.
User service -> get update and delete user information such as user data,notifications,etc

The only consumer in the system is Bot consumer.

The bot service produce a message with the relevant trading data to RabbitMQ and the consumer create a new in a go routine.
The bot will automatically buy and sell cryptocurrency using SMA and stairs as a trading strategy. 
It will produce notifications to the user about its process (buy/sell).

Each service has API struct that will implement the Make method that will later call the service.

To call any of the microservices I created a service factory that will call the registry service and get a URL of a fresh service endpoint with a success rate that is higher than 50% (some random number at the time)

Microservice call example:
```
svc := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
```

The infrastructure is consumed by the Web Application.
Even though it contains Client side and server side together it acts as API Gateway to the system.

The client side of the web application (Angularjs) is authenticating using Google auth2.

At the moment there's no implementation of bot creation through the client side but the bot service is capable of creating it so the testing file is the only way to produce.
