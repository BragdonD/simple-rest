package simplerest

// WithHttps is an option function that sets the server to use
// https over http.
func WithHttps(cert string, key string) sopts {
	return func(a any) {
		s, ok := a.(*Server)
		if ok {
			s.TLS.Key = key
			s.TLS.Cert = cert
		}
	}
}

// WithMtls is an option function that sets the server to check
// the client certificate.
func WithMtls(cert string, key string, ca []string) sopts {
	return func(a any) {
		s, ok := a.(*Server)
		if ok {
			s.TLS.Key = key
			s.TLS.Cert = cert
			s.TLS.ClientCA = ca
		}
	}
}
