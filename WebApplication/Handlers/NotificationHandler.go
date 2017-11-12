package Handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"local/UbotTrade/API/UserServiceAPI"
	"local/UbotTrade/Global"
	"local/UbotTrade/UserService/Service"
)

type NotificationsAPIHandler struct {
}

type NotificationsHandler interface {
	GetNotifications(w http.ResponseWriter, req *http.Request, ctx context.Context)
}

func NotificationsHandlerMake() NotificationsHandler {
	return &NotificationsAPIHandler{}
}

type NotificationRequestModel struct {
	ReadAll bool `json:"read_all"`
}
// get the notifications of all the bots of the logged in user
func (notificationHandler *NotificationsAPIHandler) GetNotifications(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	var requestModel NotificationRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)

	//the app will use a testing id so other people will be able to see the bots
	userId := "118103040085940455572"

	//userId := ctx.Value("UserId").(string)
	userSvc := ctx.Value(reflect.TypeOf(UserServiceAPI.UserServiceAPI{})).(*UserServiceAPI.UserServiceAPI)

	if(requestModel.ReadAll){
		data := make(map[string]interface{})
		data["IsActive"] = false
		response, err := userSvc.SetUser(Global.SetUserRequest{
			DataType:   UserService.NotificationsDataType,
			UserId:     userId,
			Data:data,
			Operation:UserService.UpdateOperation,
		})
		if (err != nil) {
			w.WriteHeader(response.Base.Status)
		}
	}

	response, err := userSvc.GetUser(Global.GetUserRequest{
		DataType:   UserService.NotificationsDataType,
		UserId:     userId,
		ActiveOnly: true,
	})

	if (err != nil) {
		w.WriteHeader(response.Base.Status)
	}

	json.NewEncoder(w).Encode(response.Data)

	//authSvc := ServiceAPIFactory.GetServiceInstance(reflect.TypeOf(AuthenticationServiceAPI.AuthenticationServiceAPI{}), registrationPath).(*AuthenticationServiceAPI.AuthenticationServiceAPI)
}
