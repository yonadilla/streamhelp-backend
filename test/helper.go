package test

import (
	"streamhelper-backend/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := DB.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		Log.Fatalf("Failed clear user data : %+v", err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := DB.First(user).Error
	assert.Nil(t, err)
	return user
}