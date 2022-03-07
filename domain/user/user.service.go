package user

import (
	"github.com/jmoiron/sqlx"
)

type userService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *userService {
	return &userService{
		db: db,
	}
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
