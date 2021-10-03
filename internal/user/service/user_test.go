package user

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/millbj92/nuboverflow-users/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Test get user by ID", func(t *testing.T){
		userStoreMock := NewMockService(mockCtrl)
		id := 1
		userStoreMock.
		EXPECT().
		GetUserByID(id).
		Return(user.User{
			ID: id,
		}, nil)

		userService := NewService(userStoreMock)
		user, err := userService.GetUserByID(
			id,
		)
		assert.NoError(t, err);
		assert.Equal(t, 1, user.ID)
	})

	t.Run("Tests inserting a user", func(t *testing.T){
		userStoreMock := NewMockService(mockCtrl)
		id := 1
		userStoreMock.
		EXPECT().
		CreateUser(&user.User{
			ID: id,
		}).
		Return(&user.User{
			ID: id,
		}, nil)

		userService := NewService(userStoreMock)
		user, err := userService.CreateUser(
			&user.User{
				ID: id,
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
	})

		t.Run("tests delete user", func(t *testing.T) {
		userStoreMock := NewMockService(mockCtrl)
		id := 1
		userStoreMock.
			EXPECT().
			DeleteUser(id).
			Return(nil)

		userService := NewService(userStoreMock)
		err := userService.DeleteUser(id)
		assert.NoError(t, err)
	})
}