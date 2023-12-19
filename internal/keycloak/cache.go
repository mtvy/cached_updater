package keycloak

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mtvy/cached_updater/internal/metrics"
	"github.com/mtvy/cached_updater/internal/userdata"
	"github.com/mtvy/cached_updater/pkg"
)

var (
	errNoCachedUser = errors.New("no cached user")
)

const (
	getCacheUserByEmail = "cache_get_user_by_email"
	getCacheUsersErr    = "cache_error_get_user_by_email"
	getUsers            = "get_users"
	getUsersErr         = "error_get_user_by_email"

	getCacheUserByID    = "cache_get_user_by_id"
	getCacheUserByIDErr = "cache_error_get_user_by_id"
	getUserByID         = "get_user_by_id"
	getUserByIDErr      = "error_get_user_by_id"
)

// cachedUser - Запись о пользователе с deadline
type cachedUser struct {
	user     *userdata.User
	deadline time.Time
}

// userCache - Хранит данные пол user'ам
// есть deadline у каждой записи user'а
type userCache struct {
	// Время жизни cachedUser
	ttl time.Duration
	// Чтобы при чтении не было проблем
	sync.RWMutex
	// Ключ - userID, значение - user с датой очистки.
	// userIDMap имеет соответсвие с элементом emailMap
	userIDMap map[string]*cachedUser
	// Ключ - email, значение - user с датой очистки.
	// emailMap имеет соответсвие с элементом userIDMap
	emailMap map[string]*cachedUser
}

func NewUserCache(ttl time.Duration, kcr pkg.UserAdapter) *userCache {
	return &userCache{
		ttl:       ttl,
		userIDMap: make(map[string]*cachedUser),
		emailMap:  make(map[string]*cachedUser),
	}
}

// SetUser - Сеттим новое значение в маппу с lock
func (c *userCache) SetUser(ctx context.Context, userID, email string, newUser userdata.User) {
	// Заводим кэшированного пользователя, который будет и в userIDMap и emailMap
	cached := cachedUser{
		user:     &newUser,
		deadline: time.Now().UTC().Add(c.ttl),
	}
	c.Lock()
	defer c.Unlock()
	c.userIDMap[userID] = &cached
	c.emailMap[email] = &cached
}

// SetUGetUserByUserIDser - Безопасно достаём User'а по userID
func (c *userCache) GetUserByUserID(ctx context.Context, userID string) (userdata.User, error) {
	c.RLock()
	defer c.RUnlock()
	// Проверяем наличие валидной записи в userIDMap
	if cached, ok := c.userIDMap[userID]; ok && cached.deadline.After(time.Now().UTC()) {
		metrics.IncKeycloakCacheEvent(getCacheUserByID)
		return *cached.user, nil
	}
	metrics.IncKeycloakCacheEvent(getCacheUserByEmail)
	return userdata.User{}, errNoCachedUser
}

// GetUserByEmail - Безопасно достаём User'а по email
func (c *userCache) GetUserByEmail(ctx context.Context, email string) (userdata.User, error) {
	c.RLock()
	defer c.RUnlock()
	if cached, ok := c.emailMap[email]; ok && cached.deadline.After(time.Now().UTC()) {
		metrics.IncKeycloakCacheEvent(getCacheUserByEmail)
		return *cached.user, nil
	}
	metrics.IncKeycloakCacheEvent(getUserByIDErr)
	return userdata.User{}, errNoCachedUser
}
