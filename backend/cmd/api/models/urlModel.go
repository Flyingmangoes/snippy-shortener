package model

import "time"

type Url struct {
    ID          string     `db:"id"`
    OriginalUrl string     `db:"originalurl"`
    ShortCode   string     `db:"shortcode"`
    IsCustom    bool       `db:"is_custom"`
    CreatedAt   time.Time  `db:"createdat"`
    ExpiryAt    time.Time  `db:"expiryat"`
}