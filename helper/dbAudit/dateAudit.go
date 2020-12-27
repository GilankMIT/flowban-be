package dbAudit

import "time"

type DateAudit struct {
	CreatedAt time.Time  `json:"-" example:"2020-08-04T14:43:10.611Z"`
	UpdatedAt time.Time  `json:"-" example:"2020-08-04T14:43:10.611Z"`
	DeletedAt *time.Time `json:"-" example:"2020-08-04T14:43:10.611Z"`
}
