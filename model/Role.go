package model

type Role struct {
	ID   int    `gorm:"primary_key"`
	Name string `json:"name"`
}
