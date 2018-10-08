package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	key      = "abcdefghijklmnopqrstuvwxyz"
	testUser = "test_user"
	testRole = "test_role"
)

func TestValidate(t *testing.T) {
	var i = 0
	for {
		if i > 10 {
			break
		}
		i++

		jwt, err := Generate(key, 2*time.Hour, testUser, testRole)
		assert.NoError(t, err)

		requester, err := Validate(key, jwt)
		assert.NoError(t, err)

		assert.Equal(t, testUser, requester.UserID)
		assert.Equal(t, testRole, requester.Role)

		t.Log(jwt)
		time.Sleep(1 * time.Second)
	}
}
