package model

import (
	"context"
	"database/sql"
	"time"
	"log/slog"
)

type UrlStoreInterface interface {
	CreateUrl(ctx context.Context, originalUrl string, sc *string, isCustom bool, expires time.Time) (*Url, error) 
	GetUrlByShortcode(ctx context.Context, sc string) (*Url, error)
	GetUrlByOriginalUrl(ctx context.Context, originalUrl string) (*Url, error)
	DeleteExpiredUrl(ctx context.Context) error 
}

type UrlStore struct {
	db *sql.DB
}

func NewUrlStore(db *sql.DB) *UrlStore {
	return &UrlStore{db: db}	
}

func (store *UrlStore) CreateUrl(ctx context.Context, originalUrl string, sc *string, isCustom bool, expires time.Time) (*Url, error) {
	newUrl := &Url{}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, 
		`INSERT INTO urls (originalurl, shortcode, is_custom, expiryat)
		VALUES ($1, $2, $3, $4)
		RETURNING id, originalurl, shortcode, createdat, expiryat, is_custom`,
		originalUrl, sc, isCustom, expires,
	).Scan(&newUrl.ID, &newUrl.OriginalUrl, &newUrl.ShortCode, &newUrl.CreatedAt, &newUrl.ExpiryAt, &newUrl.IsCustom)
	
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newUrl, nil
}

func (store *UrlStore) GetUrlByShortcode(ctx context.Context, sc string) (*Url, error) {
	newUrl := &Url{}
	err := store.db.QueryRowContext(ctx,
		`SELECT id, originalurl, shortcode, is_custom, createdat, expiryat
		FROM urls
		WHERE shortcode = $1`,
		sc,
	).Scan(&newUrl.ID, &newUrl.OriginalUrl, &newUrl.ShortCode, &newUrl.IsCustom, &newUrl.CreatedAt, &newUrl.ExpiryAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return newUrl, nil
}

func (store *UrlStore) GetUrlByOriginalUrl(ctx context.Context, originalUrl string) (*Url, error) {
	newUrl := &Url{}

	err := store.db.QueryRowContext(ctx,
		`SELECT id, originalurl, shortcode, is_custom, createdat, expiryat
		FROM urls
		WHERE originalurl = $1`,
		originalUrl,
	).Scan(&newUrl.ID, &newUrl.OriginalUrl, &newUrl.ShortCode, &newUrl.CreatedAt, &newUrl.ExpiryAt, &newUrl.IsCustom)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return newUrl, nil
}

func (store *UrlStore) DeleteExpiredUrl(ctx context.Context) error {
	result, err := store.db.ExecContext(ctx,`DELETE FROM urls WHERE expiryat < NOW()`)

	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	slog.Info("expired urls cleaned up", "rows_deleted", rowsAffected)
	return nil
}

