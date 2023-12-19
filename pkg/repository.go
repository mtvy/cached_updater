package pkg

import (
	"context"

	"github.com/mtvy/cached_updater/internal/userdata"
)

type UserAdapter interface {
	// CreateUser - Проставляем значение user'а
	CreateUser(ctx context.Context, token, realm string, user userdata.User) (string, error)
	// GetUsers - Получаем значение user'ов
	GetUsers(ctx context.Context, token, realm string, params userdata.GetUsersParams) ([]*userdata.User, error)
	// GetUserByID - Получаем значение user'а из keycloak по userID
	GetUserByID(ctx context.Context, accessToken, realm, userID string) (*userdata.User, error)

	LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*userdata.JWT, error)
	SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error
	GetCredentials(ctx context.Context, token, realm, userID string) ([]*userdata.CredentialRepresentation, error)
	DeleteCredentials(ctx context.Context, token, realm, userID, credentialID string) error
	LogoutAllSessions(ctx context.Context, accessToken, realm, userID string) error
	Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*userdata.JWT, error)

	UpdateUser(ctx context.Context, token, realm string, user userdata.User) error
}
