package webservice

import (
	"time"
)

type DB struct {
	Id          int       `json:"id"`
	Datetime    time.Time `json:"datetime"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PageHtml    string    `json:"page-html"`
}
