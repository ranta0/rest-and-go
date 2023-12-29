package user

import (
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAll() ([]User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetByID(id string) (*User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Create(user *User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) Update(id string, user *User) error {
	return s.userRepo.Update(id, user)
}

func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}

// Auth related
func (s *UserService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

func (s *UserService) Login(username, password string) (*User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return &User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &User{}, err
	}

	if err != nil {
		return &User{}, err
	}

	return user, nil
}
