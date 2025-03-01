package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type UserModule struct {
	UserController IUserController
}

func NewUserModule(pool *pgxpool.Pool, mux *http.ServeMux) *UserModule {
	repo := NewUserRepo(pool)
	service := NewUserService(repo)
	controller := NewUserController(service, mux)

	return &UserModule{UserController: controller}
}
