package cache

import (
	"context"

	"github.com/mtvy/cached_updater/internal/metrics"
	"github.com/mtvy/cached_updater/internal/userdata"
	"github.com/mtvy/cached_updater/pkg"
)

const (
	getCacheUserByEmail = "cache_get_user_by_email"
	getCacheUsersErr    = "cache_error_get_user_by_email"
	getKeycloakUsers    = "keycloak_get_users"
	getKeycloakUsersErr = "keycloak_error_get_user_by_email"
)

type cachedUsersProvider interface {
	// SetUser - Сеттим новое значение в маппу с lock
	SetUser(ctx context.Context, userID, email string, newUser userdata.User)
	// SetUGetUserByUserIDser - Безопасно достаём User'а по userID
	GetUserByUserID(ctx context.Context, userID string) (userdata.User, error)
	// GetUserByEmail - Безопасно достаём User'а по email
	GetUserByEmail(ctx context.Context, email string) (userdata.User, error)
}

type UserAdapter interface {
	pkg.UserAdapter
}

// Кэш декорирует поход в базу за UserData
type cacheDecorator struct {
	// Декорируемый интерфейс
	userAdapter UserAdapter
	// Провайдер кээшированных пользователей
	userProvider cachedUsersProvider
}

func NewCacheDecorator(userAdapter UserAdapter, userProvider cachedUsersProvider) *cacheDecorator {
	return &cacheDecorator{
		userAdapter:  userAdapter,
		userProvider: userProvider,
	}
}

// CreateUser - Проставляем значение user'а
func (c *cacheDecorator) CreateUser(ctx context.Context, token, realm string, user userdata.User) (string, error) {
	// Заводим нового пользователя в keycloak
	userID, err := c.userAdapter.CreateUser(ctx, token, realm, user)
	if err != nil {
		return userID, err
	}
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	c.userProvider.SetUser(ctx, userID, email, user)
	return userID, nil
}

// GetUserByID - Получаем значение user'а по userID
func (c *cacheDecorator) GetUserByID(ctx context.Context, accessToken, realm, userID string) (*userdata.User, error) {
	if user, err := c.userProvider.GetUserByUserID(ctx, userID); err == nil {
		return &user, nil
	}
	// Если нет - получаем и сеттим в cache
	newUserPtr, err := c.userAdapter.GetUserByID(ctx, accessToken, realm, userID)
	if err != nil {
		return newUserPtr, err
	}
	email := ""
	if newUserPtr.Email != nil {
		email = *newUserPtr.Email
	}
	c.userProvider.SetUser(ctx, userID, email, *newUserPtr)
	return newUserPtr, nil
}

// isGetUserByEmail - Если приходит только запрос на получение пользователя только по email - вернём true
func isGetUserByEmail(ctx context.Context, params userdata.GetUsersParams) bool {
	// Если нет Email в запросе - сразу вернём false
	if params.Email == nil {
		return false
	}
	// Проверяем, что не пришло других полей для отправки true
	return params.BriefRepresentation == nil &&
		params.EmailVerified == nil &&
		params.Enabled == nil &&
		params.Exact == nil &&
		params.First == nil &&
		params.FirstName == nil &&
		params.IDPAlias == nil &&
		params.IDPUserID == nil &&
		params.LastName == nil &&
		params.Max == nil &&
		params.Q == nil &&
		params.Search == nil &&
		params.Username == nil
}

// GetUsers - Получаем значение user'ов из keycloak по gocloak.GetUsersParams
func (c *cacheDecorator) GetUsers(ctx context.Context, token, realm string, params userdata.GetUsersParams) ([]*userdata.User, error) {
	// Проверяем наличие валидной записи в emailMap
	// Проверяем params на наличие только поля Email (в этом случае запишем в кэш)
	if isGetUserByEmail(ctx, params) {
		if user, err := c.userProvider.GetUserByEmail(ctx, *params.Email); err == nil {
			return []*userdata.User{&user}, nil
		}
	}
	// Если нет - получаем
	users, err := c.userAdapter.GetUsers(ctx, token, realm, params)
	if err != nil {
		metrics.IncKeycloakCacheEvent(getKeycloakUsersErr)
		return users, err
	}

	// Проставляем запись в cache
	for _, user := range users {
		email := ""
		if user.Email != nil {
			email = *user.Email
		}
		userID := ""
		if user.ID != nil {
			userID = *user.ID
		}
		c.userProvider.SetUser(ctx, userID, email, *user)
	}
	metrics.IncKeycloakCacheEvent(getKeycloakUsers)
	return users, nil
}

func (c *cacheDecorator) UpdateUser(ctx context.Context, token, realm string, user userdata.User) error {
	if err := c.userAdapter.UpdateUser(ctx, token, realm, user); err != nil {
		return err
	}
	c.userProvider.SetUser(ctx, *user.ID, *user.Email, user)
	return nil
}

func (c *cacheDecorator) LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*userdata.JWT, error) {
	return c.userAdapter.LoginClient(ctx, clientID, clientSecret, realm, scopes...)
}

func (c *cacheDecorator) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error {
	return c.userAdapter.SetPassword(ctx, token, userID, realm, password, temporary)
}

func (c *cacheDecorator) GetCredentials(ctx context.Context, token, realm, userID string) ([]*userdata.CredentialRepresentation, error) {
	return c.userAdapter.GetCredentials(ctx, token, realm, userID)
}

func (c *cacheDecorator) DeleteCredentials(ctx context.Context, token, realm, userID, credentialID string) error {
	return c.userAdapter.DeleteCredentials(ctx, token, realm, userID, credentialID)
}

func (c *cacheDecorator) LogoutAllSessions(ctx context.Context, accessToken, realm, userID string) error {
	return c.userAdapter.LogoutAllSessions(ctx, accessToken, realm, userID)
}

func (c *cacheDecorator) Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*userdata.JWT, error) {
	return c.userAdapter.Login(ctx, clientID, clientSecret, realm, username, password)
}
