package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rlawnsxo131/madre-server-v2/lib"
)

var sqlxLib = lib.NewSqlxLib()

type SocialAccountReadRepository interface {
	FindOneBySocialId(socialId string) (SocialAccount, error)
}

type socialAccountReadRepository struct {
	db *sqlx.DB
}

func NewSocialAccountReadRepository(db *sqlx.DB) SocialAccountReadRepository {
	return &socialAccountReadRepository{
		db: db,
	}
}

func (r *socialAccountReadRepository) FindOneBySocialId(socialId string) (SocialAccount, error) {
	var socialAccount SocialAccount

	query := "SELECT * FROM social_account WHERE social_id = ?"
	err := r.db.QueryRowx(query, socialId).StructScan(&socialAccount)
	if err != nil {
		customError := errors.Wrap(err, "SocialAccountRepository: FindOneBySocialId error")
		err = sqlxLib.ErrNoRowsReturnRawError(err, customError)
	}

	return socialAccount, err
}
