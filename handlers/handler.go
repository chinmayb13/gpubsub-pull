package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"hdfcbank.com/gpubsub-pull/client"
	"hdfcbank.com/gpubsub-pull/helpers"
)

type RouterConfig struct {
	Service client.PubSubClient
	Logger  *zap.Logger
}

func InitRouter(ctx context.Context, config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/count", getCountHandler(ctx, config)).Methods(http.MethodGet)
	return router
}

func getCountHandler(ctx context.Context, config RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info("routed to getCountHandler....")

		headerContentTtype := r.Header.Get("Content-Type")
		if headerContentTtype != "application/json" {
			helpers.WriteResponse(w, "application/json", "Content Type is not application/json", 0, http.StatusUnsupportedMediaType)
			return
		}

		count := config.Service.GetMessageCount()

		helpers.WriteResponse(w, "application/json", "SUCCESS", count, http.StatusOK)

	}
}
