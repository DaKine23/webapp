package oauth2

import (
	"errors"
	"os"
	"strings"

	"github.bus.zalan.do/lmineiro/go-tokens/tokens"
)

var TokenNotFound = errors.New("token not found")

type TokenProvider struct {
	predefined map[string]string
	manager    *tokens.Manager
}

type GetToken func() (string, error)

func NewTokenProvider(manager *tokens.Manager) *TokenProvider {

	return &TokenProvider{
		predefined: loadPredefinedTokens(),
		manager:    manager,
	}
}

func loadPredefinedTokens() map[string]string {
	return parsePredefinedTokens(os.Getenv("OAUTH2_ACCESS_TOKENS"))
}

func parsePredefinedTokens(env string) map[string]string {
	predefined := make(map[string]string)
	for _, pair := range strings.Split(env, ",") {
		t := strings.Split(pair, "=")
		if len(t) != 2 {
			continue
		}

		predefined[t[0]] = t[1]
	}
	return predefined
}

func (p *TokenProvider) get(name string) (string, error) {
	token, found := p.predefined[name]
	if found {
		return token, nil
	}

	if p.manager != nil {
		accessToken, err := p.manager.Get(name)
		if err != nil {
			return "", err
		}

		return accessToken.Token, nil
	}
	return "", TokenNotFound
}

func (p *TokenProvider) Tokens(name string) GetToken {
	return func() (string, error) {
		return p.get(name)
	}
}
