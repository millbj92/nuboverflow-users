package user

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/millbj92/nuboverflow-users/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Tests get all users", func(t *testing.T) {
		userStoreMock := NewMockService(mockCtrl)
		userStoreMock.
			EXPECT().
			GetAllUsers().
			Return([]user.User{}, nil)

		userService := NewService(userStoreMock)
		users, err := userService.GetAllUsers()
		assert.NoError(t, err)
		assert.IsType(t, users, []user.User{})
	})

	t.Run("Tests update user", func(t *testing.T) {

		usr := user.User{
			ID:    4,
			Email: "test@test.com",
		}
		userStoreMock := NewMockService(mockCtrl)
		userStoreMock.
			EXPECT().
			UpdateUser(usr).
			Return(user.User{
				ID:    4,
				Email: "test@test.com",
			}, nil)

		userService := NewService(userStoreMock)
		user, err := userService.UpdateUser(usr)
		assert.NoError(t, err)
		assert.Equal(t, 4, user.ID)
	})

	t.Run("Tests get user by ID", func(t *testing.T) {
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
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("Tests get user by email", func(t *testing.T) {
		userStoreMock := NewMockService(mockCtrl)
		id := 1
		email := "test@test.com"

		usr := user.User{
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			ID:         id,
			UserName:   "",
			Email:      email,
			Github:     "",
			Linkedin:   "",
			UserScore:  0,
			Bio:        "",
			Profession: "",
			WorkPlace:  "",
		}

		userStoreMock.
			EXPECT().
			GetUserByEmail(email).
			Return(usr, nil)

		userService := NewService(userStoreMock)
		user, err := userService.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, "test@test.com", user.Email)

	})

	t.Run("Tests inserting a user", func(t *testing.T) {
		userStoreMock := NewMockService(mockCtrl)
		id := 1
		email := "test@test.com"

		usr := user.User{
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			ID:         id,
			UserName:   "",
			Email:      email,
			Github:     "",
			Linkedin:   "",
			UserScore:  0,
			Bio:        "",
			Profession: "",
			WorkPlace:  "",
		}

		userStoreMock.
			EXPECT().
			CreateUser(&usr).
			Return(&usr, nil)

		userStoreMock.
			EXPECT().
			GetUserByEmail(usr.Email).
			Times(1)

		userService := NewService(userStoreMock)
		user, err := userService.CreateUser(&usr)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "test@test.com", user.Email)
	})

	t.Run("Tests delete user", func(t *testing.T) {
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
