package handler

import (
	model "backend/cmd/api/models"
	"backend/cmd/config"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type Handler struct {
	store *model.UrlStore
	cfg *config.Application
}

func NewHandler(store *model.UrlStore,cfg *config.Application) *Handler {
	return &Handler{
		store: store,
		cfg: cfg,
	}
}

func (h *Handler)ShorteningUrlHandler (w http.ResponseWriter, r *http.Request) {
	var body struct {
        Url     string  `json:"url"`
        Alias   string  `json:"alias"`
        Expires int  	`json:"expires"`
    }

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Info("[ERROR]", "Decode err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if len(body.Alias) > h.cfg.AppConf.CustomAliasLength {
		http.Error(w, "Alias had too many character", http.StatusBadRequest)
		return
	}

	expVal := body.Expires + time.Now().Day()
	slog.Info("[DEBUG]", "expires day", expVal)

	expDate := time.Date(time.Now().Year(), time.Now().Month(), expVal, 0, 0, 0, 0, time.UTC)
	slog.Info("[DEBUG]", "expire date", expDate)

	sc := ShorteningUrl(body.Url, body.Alias, expDate)

	url, err := h.store.CreateUrl(r.Context(), sc.originalurl, &sc.processedUrl, sc.isCustom, sc.expiresDate)
	if err != nil {
		slog.Info("[ERROR]", "Create err", err)
		var pgErr *pq.Error
    	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        	http.Error(w, "Alias already taken", http.StatusConflict) // 409
        	return
    	}
		http.Error(w, "Error while creating url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
 	if err = json.NewEncoder(w).Encode(url.ShortCode); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	slog.Info("======= [END MESSAGE] =======")
}

func (h *Handler)RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortcode")
	 if shortCode == "" {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

	url, err := h.store.GetUrlByShortcode(r.Context(), shortCode)
	if err != nil {
		slog.Info("[ERROR]", "err", err)
    	if errors.Is(err, sql.ErrNoRows) {	
        	http.Error(w, "Row not exist", http.StatusNotFound) 
        	return
    	}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if url == nil {
		slog.Info("[ERROR]", "err", err)
		http.Error(w, "Error while getting url", http.StatusInternalServerError)
		return
	}

	if time.Now().After(url.ExpiryAt) {
		err = h.store.DeleteExpiredUrl(r.Context())
		if err != nil {
			slog.Info("[ERROR]", "err", err)
			http.Error(w, "Error while deleting expired url", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Link has expired", http.StatusGone) 
        return
	}
	slog.Info("[PROCESS]", "redirecting", url.OriginalUrl)
	http.Redirect(w, r, url.OriginalUrl, http.StatusPermanentRedirect)
	slog.Info("======= [END MESSAGE] =======")
}

func (h *Handler) StartExpiryCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C :
			err := h.store.DeleteExpiredUrl(ctx)
			if err != nil {
				slog.Info("[ERROR] cleanup failed", "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
