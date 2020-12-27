package model

import "flowban/helper/dbAudit"

type Article struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ThumbnailUrl string `json:"thumbnail_url"`
	dbAudit.DateAudit
	dbAudit.UserAudit
}
