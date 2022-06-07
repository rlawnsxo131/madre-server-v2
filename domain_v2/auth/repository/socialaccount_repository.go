package repository

import (
	"github.com/pkg/errors"
	"github.com/rlawnsxo131/madre-server-v2/database"
	"github.com/rlawnsxo131/madre-server-v2/domain_v2/auth"
	"github.com/rlawnsxo131/madre-server-v2/utils"
)

type socialAccountRepository struct {
	db     database.Database
	mapper socialAccountEntityMapper
}

func NewSocialAccountRepository(db database.Database) auth.SocialAccountRepository {
	return &socialAccountRepository{
		db:     db,
		mapper: socialAccountEntityMapper{},
	}
}

func (r *socialAccountRepository) Create(sa *auth.SocialAccount) (string, error) {
	var id string

	query := "INSERT INTO social_account(user_id, provider, social_id)" +
		" VALUES(:user_id, :provider, :social_id)" +
		" RETURNING id"

	err := r.db.PrepareNamedGet(
		&id,
		query,
		r.mapper.toModel(sa),
	)
	if err != nil {
		return "", errors.Wrap(err, "socialaccount WriteRepository create")
	}

	return id, err
}

func (r *socialAccountRepository) FindOneBySocialIdAndProvider(socialId, provider string) (*auth.SocialAccount, error) {
	var sa auth.SocialAccount

	query := "SELECT * FROM social_account" +
		" WHERE social_id = $1" +
		" AND provider = $2"

	err := r.db.QueryRowx(query, socialId, provider).StructScan(&sa)
	if err != nil {
		customError := errors.Wrap(err, "socialaccount ReadRepository FindOneBySocialId")
		err = utils.ErrNoRowsReturnRawError(err, customError)
	}

	return r.mapper.toEntity(&sa), err
}
