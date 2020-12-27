package userDTO

type ReqAddNewUser struct {
	Password  string `json:"-"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      []int  `json:"role" binding:"required"`
}

type ReqRemoveUserByID struct {
	UserID int `json:"user_id"`
}

type ResRemoveUserByID struct {
	Message string `json:"message" example:"user successfully deleted"`
}

type ResUpdateUserByID struct {
	Message string `json:"message" example:"user successfully modified"`
}

type ReqUpdateOwnUser struct {
	Password  string `json:"password" example:"new_pass"`
	FirstName string `json:"first_name" example:"Dudung"`
	LastName  string `json:"last_name" example:"Sutoyo"`
}
