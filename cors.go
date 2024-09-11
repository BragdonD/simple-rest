package simplerest

// WithCors enables CORS management by the server
func WithCors(allowedOrigins []string, allowedHeaders []string) sopts {
	return func(a any) {
		if s, ok := a.(*Server); ok {
			s.CORS = &struct {
				AllowedOrigins []string
				AllowedHeaders []string
			}{
				AllowedOrigins: allowedOrigins,
				AllowedHeaders: allowedHeaders,
			}
		}
	}
}
