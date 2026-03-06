package server

import (
	handler "backend/cmd/api/handler"
	"backend/cmd/api/middleware"
	model "backend/cmd/api/models"
	"backend/cmd/config"
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func Start(cfg *config.Application, urlStore *model.UrlStore) {
	mux := http.NewServeMux()
	urlHandler := handler.NewHandler(urlStore, cfg)

	mux.HandleFunc("POST /api/shorten", urlHandler.ShorteningUrlHandler)
	mux.HandleFunc("GET /{shortcode}", urlHandler.RedirectUrlHandler)
	
	var h http.Handler = mux
	h = middleware.NewIPRateLimiter(
		rate.Every(time.Minute / time.Duration(cfg.RateLConf.RequestPerMinute)), 
		cfg.RateLConf.Burst,
	).RateLimit(h)
	h = middleware.CORSMiddleware(h)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go urlHandler.StartExpiryCleanup(ctx, 10*time.Minute)

	addr := net.JoinHostPort(cfg.ServConf.Host, cfg.ServConf.Port)
	slog.Info("[DEBUG]", "SERVER RUNNING AT", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}