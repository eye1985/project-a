package users

type IUserService interface {
	RegisterUser(user User) error
	GetUsers() ([]User, error)
}

type UserService struct {
	Repo IUserRepo
}

func (u *UserService) RegisterUser(user User) error {
	err := u.Repo.InsertUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUsers() ([]User, error) {
	return u.Repo.GetUsers()
}

func NewUserService(ur IUserRepo) IUserService {
	return &UserService{Repo: ur}
}

var _ IUserService = (*UserService)(nil)
