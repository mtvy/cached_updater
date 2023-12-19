package keycloak

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mtvy/cached_updater/internal/userdata"
	"github.com/stretchr/testify/require"
)

func testUserFactory(userID, email string) userdata.User {
	return userdata.User{
		ID:    GetPtr(userID),
		Email: GetPtr(email),
	}
}

func testUsersFactory(count int) []userdata.User {
	users := make([]userdata.User, count)
	for i := 0; i < count; i++ {
		users[i] = testUserFactory(strconv.Itoa(i), strconv.Itoa(i)+"@test.test")
	}
	return users
}

func TestUserCache(t *testing.T) {
	testCases := []struct {
		name    string
		user    []userdata.User
		ttl     time.Duration
		wantErr error
	}{
		{
			name:    "валидный тест",
			user:    testUsersFactory(10000),
			ttl:     time.Minute,
			wantErr: nil,
		},
		{
			name:    "тест с ошибкой при просрочке user'а в cache",
			user:    testUsersFactory(10000),
			ttl:     -time.Minute,
			wantErr: errNoCachedUser,
		},
	}

	for _, tc := range testCases {

		cache := NewUserCache(tc.ttl, nil)
		var wg sync.WaitGroup

		t.Run(tc.name+" проверяем setter user'ов в cache", func(t *testing.T) {
			wg.Add(len(tc.user))
			for _, user := range tc.user {
				go func(user userdata.User) {
					defer wg.Done()
					cache.SetUser(context.Background(), *user.ID, *user.Email, user)
				}(user)
			}
			wg.Wait()

			require.Equal(t, len(cache.userIDMap), len(tc.user))
			require.Equal(t, len(cache.emailMap), len(tc.user))
		})

		t.Run(tc.name+" проверяем GetUserByUserID и GetUserByEmail из cache", func(t *testing.T) {
			wg.Add(len(tc.user) * 2)
			for _, user := range tc.user {
				go func(user userdata.User) {
					defer wg.Done()
					cachedUser, err := cache.GetUserByUserID(context.Background(), *user.ID)
					if tc.wantErr == nil {
						require.NoError(t, err, "GetUserByUserID ошибка при получении user.ID: "+*user.ID)
						require.Equal(t, user, cachedUser, "user != cachedUser user.ID: "+*user.ID)
					} else {
						require.ErrorIs(t, tc.wantErr, err, "GetUserByUserID ожидалась ошибка при получении user.ID: "+*user.ID)
					}
				}(user)
				go func(user userdata.User) {
					defer wg.Done()
					cachedUser, err := cache.GetUserByEmail(context.Background(), *user.Email)
					if tc.wantErr == nil {
						require.NoError(t, err, "GetUserByEmail ошибка при получении user.Email: "+*user.Email)
						require.Equal(t, user, cachedUser, "user != cachedUser user.Email: "+*user.Email)
					} else {
						require.ErrorIs(t, tc.wantErr, err, "GetUserByEmail ожидалась ошибка при получении user.Email: "+*user.Email)
					}
				}(user)
			}
			wg.Wait()
		})

		t.Run(tc.name+" пытаемся параллельно записать и прочитать из cache", func(t *testing.T) {
			wg.Add(len(tc.user) * 3)
			for _, user := range tc.user {
				go func(user userdata.User) {
					defer wg.Done()
					cachedUser, err := cache.GetUserByUserID(context.Background(), *user.ID)
					if tc.wantErr == nil {
						require.NoError(t, err, "GetUserByUserID ошибка при получении user.ID: "+*user.ID)
						require.Equal(t, user, cachedUser, "user != cachedUser user.ID: "+*user.ID)
					} else {
						require.ErrorIs(t, tc.wantErr, err, "GetUserByUserID ожидалась ошибка при получении user.ID: "+*user.ID)
					}
				}(user)
				go func(user userdata.User) {
					defer wg.Done()
					cache.SetUser(context.Background(), *user.ID, *user.Email, user)
				}(user)
				go func(user userdata.User) {
					defer wg.Done()
					cachedUser, err := cache.GetUserByEmail(context.Background(), *user.Email)
					if tc.wantErr == nil {
						require.NoError(t, err, "GetUserByEmail ошибка при получении user.Email: "+*user.Email)
						require.Equal(t, user, cachedUser, "user != cachedUser user.Email: "+*user.Email)
					} else {
						require.ErrorIs(t, tc.wantErr, err, "GetUserByEmail ожидалась ошибка при получении user.Email: "+*user.Email)
					}
				}(user)
			}
			wg.Wait()
		})
	}
}
