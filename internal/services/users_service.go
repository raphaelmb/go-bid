package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphaelmb/go-bid/internal/store/pgstore"
	"golang.org/x/crypto/bcrypt"
)

const UNIQUE_CONSTRAINT_VIOLATION = "23505"

var ErrDuplicatedEmailOrPassword = errors.New("username or email already exists")

type UserService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{pool: pool, queries: pgstore.New(pool)}
}

func (us *UserService) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		PasswordHash: hash,
		Email:        email,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_CONSTRAINT_VIOLATION {
			return uuid.UUID{}, ErrDuplicatedEmailOrPassword
		}
		return uuid.UUID{}, err
	}

	return id, nil
}
