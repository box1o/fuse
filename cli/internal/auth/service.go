package auth

type Service struct {
	gateway     Gateway
	credentials CredentialStore
	browser     Browser
	waiter      Waiter
}

func NewService(gateway Gateway, credentials CredentialStore, browser Browser, waiter Waiter) *Service {
	return &Service{
		gateway:     gateway,
		credentials: credentials,
		browser:     browser,
		waiter:      waiter,
	}
}
