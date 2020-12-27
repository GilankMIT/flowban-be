package dbAudit

//UserAudit is a helper to create meta informations
//for the done of action in terms of User
type UserAudit struct {
	CreatedBy int `json:"-"`
	UpdatedBy int `json:"-"`
	DeletedBy int `json:"-"`
}
