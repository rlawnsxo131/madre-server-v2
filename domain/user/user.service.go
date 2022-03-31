package user

import (
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	UserReadRepository
	UserWriteRepository
}

type userService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Create(u User) (int64, error) {
	userWriteRepo := NewUserWriteRepository(s.db)
	lastInsertId, err := userWriteRepo.Create(u)
	return lastInsertId, err
}

func (s *userService) FindOneById(id int64) (User, error) {
	userReadRepo := NewUserReadRepository(s.db)
	user, err := userReadRepo.FindOneById(id)
	return user, err
}

func (s *userService) FindOneByUUID(uuid string) (User, error) {
	userReadRepo := NewUserReadRepository(s.db)
	user, err := userReadRepo.FindOneByUUID(uuid)
	return user, err
}
