package auth

import (
	"strings"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/services"
)

func RegisterUserStoreService() {
	err := services.AddSingleton(func() identity.UserStore {
		return &userStore{}
	})
	if err != nil {
		panic(err)
	}
}

var users = map[int]identity.User{
	1: identity.NewBasicUser(1, "albatros", "Administrators", "ToolsUsers"),
	2: identity.NewBasicUser(2, "martin", "ToolsUsers"),
}

type userStore struct{}

func (store *userStore) GetUserByID(id int) (identity.User, bool) {
	user, found := users[id]
	return user, found
}

func (store *userStore) GetUserByName(name string) (identity.User, bool) {
	for _, user := range users {
		if strings.EqualFold(user.GetDisplayName(), name) {
			return user, true
		}
	}
	return nil, false
}
