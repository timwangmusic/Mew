package config

import "golang.org/x/oauth2"

// CachedTokenSource TBD
type CachedTokenSource struct {
	RawToken oauth2.Token
}

// Token implements TokenSource
func (p *CachedTokenSource) Token() (*oauth2.Token, error) {
	return &p.RawToken, nil
}
