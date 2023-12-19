package keycloak

import "github.com/Nerzal/gocloak/v13"

type repository struct {
	*gocloak.GoCloak
}

func NewRepository(keycloakClient *gocloak.GoCloak) *repository {
	return &repository{keycloakClient}
}
