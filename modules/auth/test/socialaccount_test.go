package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rlawnsxo131/madre-server-v2/modules/auth"
	"github.com/stretchr/testify/assert"
)

func Test_SocialAccount_IsExist_IsTrue(t *testing.T) {
	assert := assert.New(t)

	sa := &auth.SocialAccount{
		ID: uuid.NewString(),
	}
	exist, err := sa.IsExist(nil)

	assert.Nil(err)
	assert.True(exist)
}
