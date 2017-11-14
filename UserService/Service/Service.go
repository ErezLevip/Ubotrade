package UserService

import (
	"log"
	"errors"
	"time"
	"reflect"
	"context"

	"gopkg.in/mgo.v2/bson"
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/API/RegistryServiceAPI"
	"github.com/erezlevip/Ubotrade/DataHandlers/MongoDB"
	"github.com/erezlevip/Ubotrade/DataHandlers/Redis"
	"github.com/erezlevip/Ubotrade/Global"
)

const NotificationsDataType = "Notifications"
const GeneralDataType = "General"
const PermissionsDataType = "Permissions"

type IUserService interface {
	GetUser(ctx context.Context, userId string, dataType string, activeOnly bool) (data []map[string]interface{}, err error)
	CreateUser(ctx context.Context, userId string, firstName string, lastName string, email string) (data map[string]interface{}, err error)
	SetUser(ctx context.Context, userId string, dataType string, operation string, data map[string]interface{}) (err error)
	GetServiceMetrics() (outputMetrics map[string]interface{}, err error)
	Init(serviceInfo Global.ServiceInformation) *UserService
}
type UserService struct {
	mongoHandler MongoDB.IDbHandler
	redisHandler DataHandlers.RedisHandler
}

const UsersCollection = "UsersData"
const UserPermissionsCollection = "UsersPermissions"
const NotificationsCollection = "BotNotifications1"

const ViewOnlyPermission = "ViewOnly"
const ViewAndEditPermission = "ViewEdit"

const CreateOperation = "Create"
const UpdateOperation = "Update"
const DeleteOperation = "Delete"

var MongoDbConfiguration map[string]string // move to redis
var RedisConfiguration DataHandlers.RedisConfiguration

var MaxUserSessionDuration = time.Duration(time.Hour * 12)

func (svc UserService) Init(serviceInfo Global.ServiceInformation) (userService *UserService) {
	log.Println(time.Now(),"Starting User Service")

	userService = &UserService{}

	ctx := context.Background()
	//register the new service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	_, err := registryService.Register(ctx, Global.RegistryRequest{ServiceInformation: serviceInfo})
	if (err != nil) {
		log.Fatal(err.Error())
	}
	//	ServiceHealth.StartHealthTicker(RedisConfiguration, serviceInfo)
	return
}

func (svc UserService) GetUser(ctx context.Context, userId string, dataType string, activeOnly bool) (data []map[string]interface{}, err error) {
	switch dataType {
	case GeneralDataType:
		data, err = getUserData(ctx , userId, UsersCollection, activeOnly)
		break
	case PermissionsDataType:
		data, err = getUserData(ctx , userId, UserPermissionsCollection, activeOnly)
		break
	case NotificationsDataType:
		data, err = getUserData(ctx , userId, NotificationsCollection, activeOnly)
		break
	case "":
		data, err = getUserData(ctx , userId, "", activeOnly)
		break
	default:
		return nil, errors.New("Invalid Data Type")
	}
	return
}

func getUserData(ctx context.Context, userId string, collection string, activeOnly bool) (data []map[string]interface{}, err error) {
	var bsonData = bson.M{"UserId": userId}
	if (activeOnly) {
		bsonData = bson.M{"UserId": userId, "IsActive": true}
	}

	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	data, err = mongoHandler.FindFMany(collection, bsonData)
	return data, err
}

func setUserData(ctx context.Context, userId string, collection string, operation string, data map[string]interface{}) (err error) {

	var bsonData = bson.M{}
	for k, v := range data {
		bsonData[k] = v
	}

	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	switch operation {

	case UpdateOperation:
		err = mongoHandler.UpdateMany(collection, bson.M{"UserId": userId}, bsonData)
		break
	case DeleteOperation:
		err = mongoHandler.Delete(collection, bson.M{"UserId": userId})
		break
	case CreateOperation:
		err = mongoHandler.Insert(collection, bsonData)
		break
	default:
		err = errors.New("Operation is required")
		break
	}

	return
}

func (svc UserService) SetUser(ctx context.Context, userId string, dataType string, operation string, data map[string]interface{}) (err error) {

	switch dataType {
	case GeneralDataType:
		err = setUserData(ctx, userId, UsersCollection, operation, data)
		break
	case PermissionsDataType:
		err = setUserData(ctx, userId, UserPermissionsCollection, operation, data)
		break
	case NotificationsDataType:
		err = setUserData(ctx, userId, NotificationsCollection, operation, data)
		break
	default:
		errors.New("Invalid Data Type")
	}
	return
}

func (svc UserService) CreateUser(ctx context.Context, userId string, firstName string, lastName string, email string) (data map[string]interface{}, err error) {

	userDataBson := bson.M{"UserId": userId, "FirstName": firstName, "LastName": lastName, "Email": email, "IsActive": true}
	err = setUserData(ctx, userId, UsersCollection, CreateOperation, userDataBson)
	if (err != nil) {
		return
	}
	userPermissionsBson := bson.M{"UserId": userId, "Permission": ViewOnlyPermission, "IsActive": true}
	err = setUserData(ctx, userId, UserPermissionsCollection, CreateOperation, userPermissionsBson)

	data = MongoDB.ConvertBsonMToMap(userDataBson)
	permissionsData := MongoDB.ConvertBsonMToMap(userPermissionsBson)

	for k, v := range permissionsData {
		data[k] = v
	}
	return
}

func (svc UserService) GetServiceMetrics() (outputMetrics map[string]interface{}, err error) {
	return make(map[string]interface{}), nil
}
