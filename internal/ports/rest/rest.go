package rest

import (
	"github.com/NikitosnikN/balance-api/internal/app"
	"github.com/NikitosnikN/balance-api/internal/app/query"
	"github.com/NikitosnikN/balance-api/internal/common/metrics"
	middleware2 "github.com/NikitosnikN/balance-api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
)

func Handler(app *app.Application) http.Handler {
	// metrics
	mHandler, mComponent := metrics.NewMetrics()

	router := chi.NewRouter()
	// Middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware2.MetricsMiddleware(mComponent))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS", "HEAD"},
		MaxAge:         300,
	}))
	//Routes
	router.Get(`/ht`, HealthcheckHandler(app.Queries.IsPoolAlive))
	router.Get(`/`, GetBalanceHandler(app.Queries.FetchBalance))
	router.Get(`/balance`, GetBalanceHandler(app.Queries.FetchBalance))
	router.Get(`/metrics`, mHandler)
	return router
}

func HealthcheckHandler(handler query.IsPoolAliveHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result, err := handler.Handle(ctx, query.IsPoolAliveQuery{})
		body := "OK"
		status := http.StatusOK

		if err != nil {
			status = http.StatusInternalServerError
			body = "Internal server error"
		}

		if result == nil || !result.Alive {
			status = http.StatusInternalServerError
			body = "Internal server error"
		}

		w.WriteHeader(status)
		render.PlainText(w, r, body)
	}
}

func GetBalanceHandler(handler query.FetchBalanceHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		address := r.URL.Query().Get("address")

		if address == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.PlainText(w, r, "address cannot be empty")
			return
		}

		blockTag := r.URL.Query().Get("blockTag")
		if blockTag == "" {
			blockTag = "latest"
		}

		balance, err := handler.Handle(ctx, query.FetchBalanceQuery{Address: address, BlockTag: blockTag})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.PlainText(w, r, "Internal server error")
		} else {
			render.JSON(w, r, balance)
		}
	}
}
