package userdata

type User struct {
	ID                         *string
	CreatedTimestamp           *int64
	Username                   *string
	Enabled                    *bool
	Totp                       *bool
	EmailVerified              *bool
	FirstName                  *string
	LastName                   *string
	Email                      *string
	FederationLink             *string
	Attributes                 *map[string][]string
	DisableableCredentialTypes *[]interface{}
	RequiredActions            *[]string
	Access                     *map[string]bool
	ClientRoles                *map[string][]string
	RealmRoles                 *[]string
	Groups                     *[]string
	ServiceAccountClientID     *string
	Credentials                *[]CredentialRepresentation
}

func (user *User) SetSiteClientID(siteClientID string) {
	if user.Attributes == nil {
		attr := make(map[string][]string, 0)
		user.Attributes = &attr
	}
	attrs := *user.Attributes
	attrs["site_client_id"] = append(attrs["site_client_id"], siteClientID)
	user.Attributes = &attrs
}

func (user *User) SetINN(inn string) {
	if user.Attributes == nil {
		attr := make(map[string][]string, 0)
		user.Attributes = &attr
	}
	attrs := *user.Attributes
	attrs["inn"] = append(attrs["inn"], inn)
	user.Attributes = &attrs
}

func (user *User) GetINN() *string {
	if user.Attributes == nil {
		return nil
	}
	attrs := *user.Attributes
	val, ok := attrs["inn"]
	if !ok || len(attrs["inn"]) == 0 {
		return nil
	}
	return &val[0]
}

func (user *User) GetPhone() *string {
	if user.Attributes == nil {
		return nil
	}
	attrs := *(user.Attributes)
	val, ok := attrs["phone"]
	if !ok || len(attrs["phone"]) == 0 {
		return nil
	}
	return &val[0]
}

type JWT struct {
	AccessToken      string
	IDToken          string
	ExpiresIn        int
	RefreshExpiresIn int
	RefreshToken     string
	TokenType        string
	NotBeforePolicy  int
	SessionState     string
	Scope            string
}

type GetUsersParams struct {
	BriefRepresentation *bool
	Email               *string
	EmailVerified       *bool
	Enabled             *bool
	Exact               *bool
	First               *int
	FirstName           *string
	IDPAlias            *string
	IDPUserID           *string
	LastName            *string
	Max                 *int
	Q                   *string
	Search              *string
	Username            *string
}

type MultiValuedHashMap struct {
	Empty      *bool    `json:"empty,omitempty"`
	LoadFactor *float32 `json:"loadFactor,omitempty"`
	Threshold  *int32   `json:"threshold,omitempty"`
}

type CredentialRepresentation struct {
	// Common part
	CreatedDate *int64  `json:"createdDate,omitempty"`
	Temporary   *bool   `json:"temporary,omitempty"`
	Type        *string `json:"type,omitempty"`
	Value       *string `json:"value,omitempty"`

	// <= v7
	Algorithm         *string             `json:"algorithm,omitempty"`
	Config            *MultiValuedHashMap `json:"config,omitempty"`
	Counter           *int32              `json:"counter,omitempty"`
	Device            *string             `json:"device,omitempty"`
	Digits            *int32              `json:"digits,omitempty"`
	HashIterations    *int32              `json:"hashIterations,omitempty"`
	HashedSaltedValue *string             `json:"hashedSaltedValue,omitempty"`
	Period            *int32              `json:"period,omitempty"`
	Salt              *string             `json:"salt,omitempty"`

	// >= v8
	CredentialData *string `json:"credentialData,omitempty"`
	ID             *string `json:"id,omitempty"`
	Priority       *int32  `json:"priority,omitempty"`
	SecretData     *string `json:"secretData,omitempty"`
	UserLabel      *string `json:"userLabel,omitempty"`
}
