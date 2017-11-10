package main

import (
	"local/UbotTrade/WebApplication/Handlers"
	"local/UbotTrade/WebApplication/Middlewares"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"local/UbotTrade/Logger"
)

func main() {
	r := mux.NewRouter()
	botHandler := Handlers.BotInfoHandlerMake()
	authHandler := Handlers.AuthHandlerMake()
	notificationsHandler := Handlers.NotificationsHandlerMake()

	Logger.SetGlobalLogger()
	var resolverMiddleware = Middlewares.ResolverMiddlewareMake()
	resolverMiddleware.Register()
	var authMiddleware = Middlewares.AuthMiddleWareMake()
	authMiddleware.Register()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public"))).Methods("GET")
	//r.HandleFunc("/GetBotActivity", botHandler.GetBotActivity)
	r.HandleFunc("/GetBotInformation", Middlewares.BeginRequest(botHandler.GetBotInformation, false)).Methods("POST")
	r.HandleFunc("/GetAllActiveBots", Middlewares.BeginRequest(botHandler.GetAllActiveBots, false)).Methods("POST")
	r.HandleFunc("/GetLastActivities", Middlewares.BeginRequest(botHandler.GetLastActivities, false)).Methods("POST")
	r.HandleFunc("/GetBotProfits", Middlewares.BeginRequest(botHandler.GetBotProfits, false)).Methods("POST")
	//r.HandleFunc("/CreateNewBot", botHandler.CreateNewBot)
	//r.HandleFunc("/CancelBot", botHandler.CancelBot)
	r.HandleFunc("/GetBotTickerData", Middlewares.BeginRequest(botHandler.GetBotTickerData, false)).Methods("POST")
	r.HandleFunc("/Login", Middlewares.BeginRequest(authHandler.Login, true)).Methods("POST")
	r.HandleFunc("/GetNotifications", Middlewares.BeginRequest(notificationsHandler.GetNotifications, false)).Methods("POST")

	port := "8000"
	log.Println(time.Now(), "WebApplication is now running on port:"+port)
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
