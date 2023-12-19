package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mtvy/cached_updater/internal/userdata"
	"github.com/mtvy/cached_updater/pkg"
)

type adapter struct {
	repo *repository
}

func NewAdapter(repo *repository) pkg.UserAdapter {
	return &adapter{repo: repo}
}

// multiValuedHashMapToKeyCloak - переводим *userdata.MultiValuedHashMap к *gocloak.MultiValuedHashMap
func multiValuedHashMapToKeyCloak(multiValuedHashMap *userdata.MultiValuedHashMap) *gocloak.MultiValuedHashMap {
	if multiValuedHashMap == nil {
		return nil
	}
	return &gocloak.MultiValuedHashMap{
		Empty:      multiValuedHashMap.Empty,
		LoadFactor: multiValuedHashMap.LoadFactor,
		Threshold:  multiValuedHashMap.Threshold,
	}
}

// credentialRepresentationToKeyCloak - переводим *[]userdata.CredentialRepresentation к *[]gocloak.CredentialRepresentation
func credentialRepresentationToKeyCloak(credentialRepresentations *[]userdata.CredentialRepresentation) *[]gocloak.CredentialRepresentation {
	if credentialRepresentations == nil {
		return nil
	}
	keyclockCRS := make([]gocloak.CredentialRepresentation, len(*credentialRepresentations))
	for i, credentialRepresentation := range *credentialRepresentations {
		keyclockCRS[i] = gocloak.CredentialRepresentation{
			CreatedDate: credentialRepresentation.CreatedDate,
			Temporary:   credentialRepresentation.Temporary,
			Type:        credentialRepresentation.Type,
			Value:       credentialRepresentation.Value,

			Algorithm:         credentialRepresentation.Algorithm,
			Config:            multiValuedHashMapToKeyCloak(credentialRepresentation.Config),
			Counter:           credentialRepresentation.Counter,
			Device:            credentialRepresentation.Device,
			Digits:            credentialRepresentation.Digits,
			HashIterations:    credentialRepresentation.HashIterations,
			HashedSaltedValue: credentialRepresentation.HashedSaltedValue,
			Period:            credentialRepresentation.Period,
			Salt:              credentialRepresentation.Salt,

			ID:         credentialRepresentation.ID,
			Priority:   credentialRepresentation.Priority,
			SecretData: credentialRepresentation.SecretData,
			UserLabel:  credentialRepresentation.UserLabel,
		}
	}
	return &keyclockCRS
}

// userToKeyCloak - переводим userdata.User к gocloak.User
func userToKeyCloak(user userdata.User) gocloak.User {
	return gocloak.User{
		ID:                         user.ID,
		CreatedTimestamp:           user.CreatedTimestamp,
		Username:                   user.Username,
		Enabled:                    user.Enabled,
		Totp:                       user.Totp,
		EmailVerified:              user.EmailVerified,
		FirstName:                  user.FirstName,
		LastName:                   user.LastName,
		Email:                      user.Email,
		FederationLink:             user.FederationLink,
		Attributes:                 user.Attributes,
		DisableableCredentialTypes: user.DisableableCredentialTypes,
		RequiredActions:            user.RequiredActions,
		Access:                     user.Access,
		ClientRoles:                user.ClientRoles,
		RealmRoles:                 user.RealmRoles,
		Groups:                     user.Groups,
		ServiceAccountClientID:     user.ServiceAccountClientID,
		Credentials:                credentialRepresentationToKeyCloak(user.Credentials),
	}
}

// getUsersParamsToKeyCloak - переводим userdata.GetUsersParams к gocloak.GetUsersParams
func getUsersParamsToKeyCloak(params userdata.GetUsersParams) gocloak.GetUsersParams {
	return gocloak.GetUsersParams{
		BriefRepresentation: params.BriefRepresentation,
		Email:               params.Email,
		EmailVerified:       params.EmailVerified,
		Enabled:             params.Enabled,
		Exact:               params.Exact,
		First:               params.First,
		FirstName:           params.FirstName,
		IDPAlias:            params.IDPAlias,
		IDPUserID:           params.IDPAlias,
		LastName:            params.LastName,
		Max:                 params.Max,
		Q:                   params.Q,
		Search:              params.Search,
		Username:            params.Username,
	}
}

// multiValuedHashMapToService - переводим *gocloak.MultiValuedHashMap к *userdata.MultiValuedHashMap
func multiValuedHashMapToService(multiValuedHashMap *gocloak.MultiValuedHashMap) *userdata.MultiValuedHashMap {
	if multiValuedHashMap == nil {
		return nil
	}
	return &userdata.MultiValuedHashMap{
		Empty:      multiValuedHashMap.Empty,
		LoadFactor: multiValuedHashMap.LoadFactor,
		Threshold:  multiValuedHashMap.Threshold,
	}
}

// ptrCredentialRepresentationToService - переводим *[]gocloak.CredentialRepresentation к *[]userdata.CredentialRepresentation
func ptrCredentialRepresentationToService(keycloakCRS *[]gocloak.CredentialRepresentation) *[]userdata.CredentialRepresentation {
	if keycloakCRS == nil {
		return nil
	}

	credentialRepresentations := make([]userdata.CredentialRepresentation, len(*keycloakCRS))
	for i, credentialRepresentation := range *keycloakCRS {
		credentialRepresentations[i] = userdata.CredentialRepresentation{
			CreatedDate: credentialRepresentation.CreatedDate,
			Temporary:   credentialRepresentation.Temporary,
			Type:        credentialRepresentation.Type,
			Value:       credentialRepresentation.Value,

			Algorithm:         credentialRepresentation.Algorithm,
			Config:            multiValuedHashMapToService(credentialRepresentation.Config),
			Counter:           credentialRepresentation.Counter,
			Device:            credentialRepresentation.Device,
			Digits:            credentialRepresentation.Digits,
			HashIterations:    credentialRepresentation.HashIterations,
			HashedSaltedValue: credentialRepresentation.HashedSaltedValue,
			Period:            credentialRepresentation.Period,
			Salt:              credentialRepresentation.Salt,

			ID:         credentialRepresentation.ID,
			Priority:   credentialRepresentation.Priority,
			SecretData: credentialRepresentation.SecretData,
			UserLabel:  credentialRepresentation.UserLabel,
		}
	}
	return &credentialRepresentations
}

// credentialRepresentationToService - переводим []*gocloak.CredentialRepresentation к []*userdata.CredentialRepresentation
func credentialRepresentationToService(keycloakCRS []*gocloak.CredentialRepresentation) []*userdata.CredentialRepresentation {
	if keycloakCRS == nil {
		return nil
	}

	credentialRepresentations := make([]*userdata.CredentialRepresentation, len(keycloakCRS))
	for i, credentialRepresentation := range keycloakCRS {
		credentialRepresentations[i] = &userdata.CredentialRepresentation{
			CreatedDate: credentialRepresentation.CreatedDate,
			Temporary:   credentialRepresentation.Temporary,
			Type:        credentialRepresentation.Type,
			Value:       credentialRepresentation.Value,

			Algorithm:         credentialRepresentation.Algorithm,
			Config:            multiValuedHashMapToService(credentialRepresentation.Config),
			Counter:           credentialRepresentation.Counter,
			Device:            credentialRepresentation.Device,
			Digits:            credentialRepresentation.Digits,
			HashIterations:    credentialRepresentation.HashIterations,
			HashedSaltedValue: credentialRepresentation.HashedSaltedValue,
			Period:            credentialRepresentation.Period,
			Salt:              credentialRepresentation.Salt,

			ID:         credentialRepresentation.ID,
			Priority:   credentialRepresentation.Priority,
			SecretData: credentialRepresentation.SecretData,
			UserLabel:  credentialRepresentation.UserLabel,
		}
	}
	return credentialRepresentations
}

// userToService - переводим *gocloak.User к *userdata.User
func userToService(keycloakUser *gocloak.User) *userdata.User {
	if keycloakUser == nil {
		return nil
	}
	return &userdata.User{
		ID:                         keycloakUser.ID,
		CreatedTimestamp:           keycloakUser.CreatedTimestamp,
		Username:                   keycloakUser.Username,
		Enabled:                    keycloakUser.Enabled,
		Totp:                       keycloakUser.Totp,
		EmailVerified:              keycloakUser.EmailVerified,
		FirstName:                  keycloakUser.FirstName,
		LastName:                   keycloakUser.LastName,
		Email:                      keycloakUser.Email,
		FederationLink:             keycloakUser.FederationLink,
		Attributes:                 keycloakUser.Attributes,
		DisableableCredentialTypes: keycloakUser.DisableableCredentialTypes,
		RequiredActions:            keycloakUser.RequiredActions,
		Access:                     keycloakUser.Access,
		ClientRoles:                keycloakUser.ClientRoles,
		RealmRoles:                 keycloakUser.RealmRoles,
		Groups:                     keycloakUser.Groups,
		ServiceAccountClientID:     keycloakUser.ServiceAccountClientID,
		Credentials:                ptrCredentialRepresentationToService(keycloakUser.Credentials),
	}
}

// usersToService - переводим []*gocloak.User к []*userdata.User
func usersToService(keycloakUsers []*gocloak.User) []*userdata.User {
	if keycloakUsers == nil {
		return nil
	}
	users := make([]*userdata.User, len(keycloakUsers))
	for i, keycloakUser := range keycloakUsers {
		users[i] = userToService(keycloakUser)
	}
	return users
}

// jwtToService - переводим *gocloak.JWT к *userdata.JWT
func jwtToService(keycloakJWT *gocloak.JWT) *userdata.JWT {
	if keycloakJWT == nil {
		return nil
	}
	return &userdata.JWT{
		AccessToken:      keycloakJWT.AccessToken,
		IDToken:          keycloakJWT.IDToken,
		ExpiresIn:        keycloakJWT.ExpiresIn,
		RefreshExpiresIn: keycloakJWT.RefreshExpiresIn,
		RefreshToken:     keycloakJWT.RefreshToken,
		TokenType:        keycloakJWT.TokenType,
		NotBeforePolicy:  keycloakJWT.NotBeforePolicy,
		SessionState:     keycloakJWT.SessionState,
		Scope:            keycloakJWT.Scope,
	}
}

// CreateUser - Проставляем значение user'а делая запрос в keycloak
func (a *adapter) CreateUser(ctx context.Context, token, realm string, user userdata.User) (string, error) {
	return a.repo.CreateUser(ctx, token, realm, userToKeyCloak(user))
}

// GetUsers - Получаем значение user'ов из keycloak
func (a *adapter) GetUsers(ctx context.Context, token, realm string, params userdata.GetUsersParams) ([]*userdata.User, error) {
	keycloakUser, err := a.repo.GetUsers(ctx, token, realm, getUsersParamsToKeyCloak(params))
	return usersToService(keycloakUser), err
}

// GetUserByID - Получаем значение user'а из keycloak по userID
func (a *adapter) GetUserByID(ctx context.Context, accessToken, realm, userID string) (*userdata.User, error) {
	keycloakUser, err := a.repo.GetUserByID(ctx, accessToken, realm, userID)
	return userToService(keycloakUser), err
}

func (a *adapter) LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*userdata.JWT, error) {
	keycloakJWT, err := a.repo.LoginClient(ctx, clientID, clientSecret, realm, scopes...)
	return jwtToService(keycloakJWT), err
}

func (a *adapter) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error {
	return a.repo.SetPassword(ctx, token, userID, realm, password, temporary)
}

func (a *adapter) GetCredentials(ctx context.Context, token, realm, userID string) ([]*userdata.CredentialRepresentation, error) {
	keycloakCR, err := a.repo.GetCredentials(ctx, token, realm, userID)
	return credentialRepresentationToService(keycloakCR), err
}

func (a *adapter) DeleteCredentials(ctx context.Context, token, realm, userID, credentialID string) error {
	return a.repo.DeleteCredentials(ctx, token, realm, userID, credentialID)
}

func (a *adapter) LogoutAllSessions(ctx context.Context, accessToken, realm, userID string) error {
	return a.repo.LogoutAllSessions(ctx, accessToken, realm, userID)
}

func (a *adapter) Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*userdata.JWT, error) {
	keycloakJWT, err := a.repo.Login(ctx, clientID, clientSecret, realm, username, password)
	return jwtToService(keycloakJWT), err
}

func (a *adapter) UpdateUser(ctx context.Context, token, realm string, user userdata.User) error {
	keycloakUser := userToKeyCloak(user)
	return a.repo.UpdateUser(ctx, token, realm, keycloakUser)
}
