package validators

type User struct {
	Username string  `json:"username" validate:"required,min=3,max=150"`
	About    *string `json:"about" validate:"omitempty,max=256"`
	Age      *uint8  `json:"age" validate:"omitempty,min=1,max=255"`
}
