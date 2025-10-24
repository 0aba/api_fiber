package validators_request

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=150"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type CommonDataUpdateRequest struct {
	Username        string  `json:"username" validate:"required,min=3,max=150"`
	CurrentPassword string  `json:"current_password" validate:"required,min=8,max=128"`
	About           *string `json:"about" validate:"omitempty,max=256"`
	Age             *uint8  `json:"age" validate:"omitempty,min=1,max=255"`
}

type PasswordUpdateRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=150"`
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=128"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type DisabledUserRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=150"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=CurrentPassword"`
}

type UsersPageRequest struct {
	SizePage uint64 `json:"size_page" validate:"required,min=1,max=65535"`
	Page     uint64 `json:"page" validate:"required,min=1,max=9223372036854775807"`
}
