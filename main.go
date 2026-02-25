package main

import (
	"context"
	"log"
	"net/http"
	"root/handlers/api"
	"root/repo"
	"root/services/login"
	"root/services/middleware"
	"root/services/sendCoins"
	"root/services/shopping"
	"root/services/userInfo"
	"root/store/postgres"
	"root/store/postgres/config"
)

func main() {
	//конфиг для БД
	cnf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't parse config, err: %v", err)
	}

	//клиент для БД
	pstgrCl, err := postgres.NewClient(context.Background(), cnf)
	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	defer pstgrCl.DbPool.Close()

	//репозиторий с запросами к БД
	repoPg := repo.NewUserRepo(pstgrCl)

	//инициализация сервисов
	lgS := login.NewLoginService(repoPg)
	scS := sendCoins.NewSendCoinsService(repoPg)
	shS := shopping.NewShoppingService(repoPg)
	uiS := userInfo.NewUsernInfoService(repoPg)

	h := api.NewHandlers(lgS, scS, shS, uiS)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth", h.Auth)
	mux.Handle("GET /api/buy/{item}", middleware.Auth(http.HandlerFunc(h.BuyItem)))
	mux.Handle("POST /api/sendCoin", middleware.Auth(http.HandlerFunc(h.SendCoins)))
	mux.Handle("GET /api/info", middleware.Auth(http.HandlerFunc(h.Info)))
	serv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	log.Print("Server is starting on 8080")
	serv.ListenAndServe()
}
