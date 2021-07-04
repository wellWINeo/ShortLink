package main

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/wellWINeo/ShortLink"
	"github.com/wellWINeo/ShortLink/pkg/handler"
	"github.com/wellWINeo/ShortLink/pkg/repository"
	"github.com/wellWINeo/ShortLink/pkg/service"
)

func main() {
	if err := initConfig(); err != nil {
		log.Printf("Can't read config: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := repository.NewMongo(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetInt("db.port"),
		Ctx: &ctx,
	})

	defer func() {
		cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Can't close db connection: %s", err.Error())
		}
	}()

	if err != nil {
		log.Fatalf("Can't connect to mongoDB: %s", err.Error())
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Can't ping db: %s", err.Error())
	}

	db := client.Database(viper.GetString("db.name"))

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services,
		viper.GetString("server.static_files"),
		viper.GetString("server.domain"))
	srv := new(ShortLink.Server)

	err = srv.Run(viper.GetInt("server.port"), handlers.InitRoutes())
	if err != nil {
		log.Fatalf("Can't start server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
