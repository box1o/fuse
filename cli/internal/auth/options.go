package auth

type LoginOptions struct {
	ClientName  string
	OpenBrowser bool
	Observer    LoginObserver
}

func (o LoginOptions) clientName() string {
	if o.ClientName == "" {
		return "Fuse CLI"
	}

	return o.ClientName
}
