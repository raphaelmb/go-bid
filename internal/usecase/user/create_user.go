package user

import (
	"context"

	"github.com/raphaelmb/go-bid/internal/validator"
)

type CreateUserReq struct {
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Bio          string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")

	return eval
}
