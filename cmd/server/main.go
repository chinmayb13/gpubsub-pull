package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"hdfcbank.com/gpubsub-pull/client"
	"hdfcbank.com/gpubsub-pull/config"
	"hdfcbank.com/gpubsub-pull/handlers"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()

	envConfig, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("error reading env file", zap.String("err", err.Error()))
	}

	ctx := context.Background()
	pubsubClient := client.GetClientAndSubscription(ctx, &client.SubscriberConfig{
		ProjectID: envConfig.PubSub.ProjectID,
		SubID:     envConfig.PubSub.SubID,
		Logger:    logger,
	})
	defer pubsubClient.CloseClient()

	router := handlers.InitRouter(ctx, handlers.RouterConfig{
		Service: pubsubClient,
		Logger:  logger,
	})
	logger.Info("starting gpubsub-pull on...", zap.String("PORT", envConfig.DeployConfig.Port))

	certPath := strings.Join([]string{envConfig.DeployConfig.AppDir, "cert"}, "/")
	go func() {
		log.Fatal(http.ListenAndServeTLS(":"+envConfig.DeployConfig.Port, certPath+"/local.crt", certPath+"/local.key", router))
	}()

	err = pubsubClient.PullMessages(ctx, envConfig.DeployConfig.NumGoRoutine)
	if err != nil {
		logger.Fatal("failed pulling messages....aborting!!!", zap.String("err", err.Error()))
	}

}
