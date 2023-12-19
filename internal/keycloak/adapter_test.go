package keycloak

import (
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mtvy/cached_updater/internal/userdata"
	"github.com/stretchr/testify/assert"
)

func GetPtr[T any](v T) *T {
	return &v
}

func Test_userToKeyCloak(t *testing.T) {
	testCases := []struct {
		name string
		user userdata.User
		want gocloak.User
	}{
		{
			name: "валидный тест с конвертированием",
			user: userdata.User{
				ID:            GetPtr("id"),
				Username:      GetPtr("XXXXXXXX"),
				FirstName:     GetPtr("firstName"),
				LastName:      GetPtr("lastName"),
				Email:         GetPtr("email"),
				Enabled:       GetPtr(true),
				EmailVerified: GetPtr(true),
				Attributes:    GetPtr(map[string][]string{"attr": {"val"}}),
			},
			want: gocloak.User{
				ID:            GetPtr("id"),
				Username:      GetPtr("XXXXXXXX"),
				FirstName:     GetPtr("firstName"),
				LastName:      GetPtr("lastName"),
				Email:         GetPtr("email"),
				Enabled:       GetPtr(true),
				EmailVerified: GetPtr(true),
				Attributes:    GetPtr(map[string][]string{"attr": {"val"}}),
			},
		},
		{
			name: "валидный тест с nil",
			user: userdata.User{
				ID:            nil,
				Username:      nil,
				FirstName:     nil,
				LastName:      nil,
				Email:         nil,
				Enabled:       nil,
				EmailVerified: nil,
				Attributes:    nil,
			},
			want: gocloak.User{
				ID:            nil,
				Username:      nil,
				FirstName:     nil,
				LastName:      nil,
				Email:         nil,
				Enabled:       nil,
				EmailVerified: nil,
				Attributes:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, userToKeyCloak(tc.user))
		})
	}
}

func Test_userToService(t *testing.T) {
	testCases := []struct {
		name string
		user gocloak.User
		want userdata.User
	}{
		{
			name: "валидный тест с конвертированием",
			user: gocloak.User{
				ID:            GetPtr("id"),
				Username:      GetPtr("XXXXXXXX"),
				FirstName:     GetPtr("firstName"),
				LastName:      GetPtr("lastName"),
				Email:         GetPtr("email"),
				Enabled:       GetPtr(true),
				EmailVerified: GetPtr(true),
				Attributes:    GetPtr(map[string][]string{"attr": {"val"}}),
			},
			want: userdata.User{
				ID:            GetPtr("id"),
				Username:      GetPtr("XXXXXXXX"),
				FirstName:     GetPtr("firstName"),
				LastName:      GetPtr("lastName"),
				Email:         GetPtr("email"),
				Enabled:       GetPtr(true),
				EmailVerified: GetPtr(true),
				Attributes:    GetPtr(map[string][]string{"attr": {"val"}}),
			},
		},
		{
			name: "валидный тест с nil",
			user: gocloak.User{
				ID:            nil,
				Username:      nil,
				FirstName:     nil,
				LastName:      nil,
				Email:         nil,
				Enabled:       nil,
				EmailVerified: nil,
				Attributes:    nil,
			},
			want: userdata.User{
				ID:            nil,
				Username:      nil,
				FirstName:     nil,
				LastName:      nil,
				Email:         nil,
				Enabled:       nil,
				EmailVerified: nil,
				Attributes:    nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, *userToService(&tc.user))
		})
	}
}
