package odata

type BaseAuthorization struct {
	Name     string
	Password string
}

func NewBaseCredentials(username, password string) *BaseAuthorization {
	return &BaseAuthorization{
		Name:     username,
		Password: password,
	}
}
