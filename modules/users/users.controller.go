package users

import (
	"encoding/json"
	"net/http"
)

type IUserController interface {
	GetUsers()
	RegisterUser()
}

type UserController struct {
	mux         *http.ServeMux
	UserService IUserService
}

func NewUserController(userService IUserService, mux *http.ServeMux) *UserController {
	return &UserController{UserService: userService, mux: mux}
}

var _ IUserController = (*UserController)(nil)

func (c *UserController) GetUsers() {
	c.mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		users, err := c.UserService.GetUsers()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (c *UserController) RegisterUser() {
	c.mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid request content-type", http.StatusUnsupportedMediaType)
			return
		}

		user := &User{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Username == "" {
			http.Error(w, "Missing required field: username", http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			http.Error(w, "Missing required field: email", http.StatusBadRequest)
			return
		}

		err = c.UserService.RegisterUser(*user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&user)
	})
}
