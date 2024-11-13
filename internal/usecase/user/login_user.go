package user

import (
	"context"

	"github.com/raphaelmb/go-bid/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lu LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(lu.Email, validator.EmailRX), "email", "must be a valid email")
	eval.CheckField(validator.NotBlank(lu.Password), "password", "this field cannot be blank")

	return eval
}
