package Handlers

import (
	"context"
	"encoding/json"
	"local/UbotTrade/API/AuthenticationServiceAPI"
	"local/UbotTrade/Global"
	"log"
	"net/http"
	"reflect"
)

type AuthorizeUserRequestModel struct {
	Uid       string `json:"uid"`
	SessionId string `json:"session_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type AuthApiHandler struct {
}

type AuthHandler interface {
	Login(w http.ResponseWriter, req *http.Request, ctx context.Context)
}

func AuthHandlerMake() AuthHandler {
	return &AuthApiHandler{}
}

func (authHandler *AuthApiHandler) Login(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	var requestModel AuthorizeUserRequestModel
	_ = json.NewDecoder(req.Body).Decode(&requestModel)

	authSvc := ctx.Value(reflect.TypeOf(AuthenticationServiceAPI.AuthenticationServiceAPI{})).(*AuthenticationServiceAPI.AuthenticationServiceAPI)

	var res = make(map[string]interface{})
	res["IsAuthorized"] = false
	if requestModel.SessionId != "" {
		authRequest := Global.LoginRequest{
			SessionId: requestModel.SessionId,
			ClientId:  requestModel.Uid,
			Email:     requestModel.Email,
			LastName:  requestModel.LastName,
			FirstName: requestModel.FirstName,
		}
		authRes, err := authSvc.Login(authRequest)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		if authRes.UserId != "" {
			res["IsAuthorized"] = true
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(w).Encode(res)

	return
}
