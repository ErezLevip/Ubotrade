package AuthenticationService

import (
	"errors"
	"local/UbotTrade/API/RegistryServiceAPI"
	"local/UbotTrade/API/ServiceAPIFactory"
	"local/UbotTrade/API/UserServiceAPI"
	"local/UbotTrade/DataHandlers/Redis"
	"local/UbotTrade/Global"
	"local/UbotTrade/UserService/Service"
	"log"
	"reflect"
	"time"
	"context"
)

type IAuthenticationService interface {
	GetToken(ctx context.Context, clientId string) (string, error)
	ValidateToken(ctx context.Context, token string) (IsValid bool, err error)
	Login(ctx context.Context, clientId string, firstName string, lastName string, email string, sessionId string) (string, error)
	GetServiceMetrics() (outputMetrics map[string]interface{}, err error)
	Init(serviceInfo Global.ServiceInformation) *AuthenticationService
}
type AuthenticationService struct {
}

const UsersCollection = "Users"
const UserPermissionsCollection = "UsersPermissions"
const ActiveSessionPrefix = "ActiveSession"

const ViewOnlyPermission = "ViewOnly"
const ViewAndEditPermission = "ViewEdit"

var MaxUserSessionDuration = time.Duration(time.Hour * 12)

func (svc AuthenticationService) Init(serviceInfo Global.ServiceInformation) (authenticationService *AuthenticationService) {

	log.Println(time.Now(),"Starting Authentication Service")

	authenticationService = &AuthenticationService{}

	ctx := context.Background()
	//register the new service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	_, err := registryService.Register(ctx, Global.RegistryRequest{ServiceInformation: serviceInfo})

	if (err != nil) {
		log.Fatal(err.Error())
	}

	return
}

func (svc AuthenticationService) Login(ctx context.Context, clientId string, firstName string, lastName string, email string, sessionId string) (userId string, err error) {

	redisHandler := ctx.Value("RedisHandler").(DataHandlers.RedisHandler)

	err = nil
	userId = clientId

	if (sessionId == "") {
		log.Println(time.Now(),"session id is empty") //sessionId can never be empty, since it is provided by the oauth providers google/facebook
		err = errors.New("UnAuthorized")
		return
	}

	var redisSessionKey = ActiveSessionPrefix + ":" + sessionId

	var sessionValue string
	sessionValue, err = redisHandler.Get(redisSessionKey)
	log.Println(time.Now(),"session value", sessionValue)
	if err != nil && sessionValue != "" {
		log.Println(err.Error())
		return
	} else if sessionValue != "" {
		if sessionValue[0] == '"' {
			sessionValue = sessionValue[1: len(sessionValue)-1]
		}
		userId = sessionValue
		return
	} else if clientId == "" { // in case there's no active session and there's no client id, there is no primary identification method
		log.Println(time.Now(),"clientid is empty")
		err = errors.New("UnAuthorized")
		return
	}

	userSvc := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(UserServiceAPI.UserServiceAPI{})).(*UserServiceAPI.UserServiceAPI)

	var userData Global.GetUserResponse
	userData, err = userSvc.GetUser(Global.GetUserRequest{
		UserId:     userId,
		ActiveOnly: true,
		DataType:   UserService.GeneralDataType,
	})

	if err != nil {
		log.Println(err.Error())
		return
	}

	if userData.Data == nil || len(userData.Data) == 0 {
		_, err = userSvc.CreateUser(Global.CreateUserRequest{
			UserId:    userId,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		var props = make(map[string]interface{})
		props["FirstName"] = firstName
		props["LastName"] = lastName
		props["Email"] = email

		userSvc.SetUser(Global.SetUserRequest{
			DataType:  UserService.NotificationsDataType,
			UserId:    userId,
			Data:      props,
			Operation: UserService.CreateOperation,
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	var permissions Global.GetUserResponse
	permissions, err = userSvc.GetUser(Global.GetUserRequest{
		UserId:     userId,
		ActiveOnly: true,
		DataType:   UserService.PermissionsDataType,
	})

	if err != nil {
		log.Println(err.Error())
		return
	}

	if userData.Data != nil && len(userData.Data) > 0 && permissions.Data != nil && len(permissions.Data) > 0 {
		userId = userData.Data[0]["UserId"].(string)
		redisHandler.Set(redisSessionKey, userId, MaxUserSessionDuration)
		return
	}
	err = errors.New("UnAuthorized")
	return
}

func (svc AuthenticationService) GetToken(ctx context.Context, clientId string) (string, error) {
	return "", nil
}

func (svc AuthenticationService) ValidateToken(ctx context.Context, token string) (IsValid bool, err error) {

	return true, nil
}

func (svc AuthenticationService) GetServiceMetrics() (outputMetrics map[string]interface{}, err error) {
	return make(map[string]interface{}), nil
}
