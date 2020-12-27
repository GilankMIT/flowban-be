package model

import "flowban/helper/dbAudit"

type User struct {
	ID                        int    `json:"id" gorm:"primary_key"`
	Password                  string `json:"-" gorm:"size:180"`
	Email                     string `json:"email" gorm:"size:150"`
	FirstName                 string `json:"first_name" gorm:"size:50"`
	LastName                  string `json:"last_name" gorm:"size:50"`
	LastActive                int64  `json:"last_active"`
	UserRole                  []Role `gorm:"many2many:user_roles;" json:"user_role"`
	EmailVerificationToken    string `json:"email_verification_token"  gorm:"size:128"`
	EmailVerificationTokenExp int64  `json:"email_verification_token_exp"`
	Address                   string `json:"address"`
	Telephone                 string `json:"telephone"`
	Position                  string `json:"position"`
	Status                    int    `json:"status"`
	dbAudit.DateAudit
	dbAudit.UserAudit
}
