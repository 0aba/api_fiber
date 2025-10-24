package validators_response

import (
	"api_fiber/src/validators"
)

type ListUsers struct {
	Users       []validators.User `json:"users" validate:"required,dive"`
	AmountPage  uint64            `json:"amount_page" validate:"required,min=0,max=9223372036854775807"`
	SizePage    uint64            `json:"size_page" validate:"required,min=0,max=65535"`
	CurrentPage uint64            `json:"current_page" validate:"required,min=0,max=9223372036854775807"`
}
