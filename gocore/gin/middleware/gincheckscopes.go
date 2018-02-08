package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.bus.zalan.do/ale/gocore/logger"
	"github.bus.zalan.do/ale/gocore/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var sc ScopeChecker = &defaultScopeChecker{}

//NoScopeChecks switch off scope checks for testing and run locally
func NoScopeChecks() {

	Init(&noOpScopeChecker{})

}

//Init is not needed for produktion!
func Init(scoche ScopeChecker) {

	sc = scoche
}

func CheckScopes(scopes ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		flowid := c.Keys[FlowIDKey].(string)

		for _, scope := range scopes {
			err := sc.checkScope(scope, c.Request.Header)
			if err != nil {
				logger.LogInfo(flowid, "%s : %v", "Scope "+scope+" is missing!", err)

				c.AbortWithError(401, fmt.Errorf("Scope "+scope+" is missing! : %v", err))
			}
		}
	}
}

type ScopeChecker interface {
	checkScope(scope string, h http.Header) error
}

type noOpScopeChecker struct{}

func (nosc *noOpScopeChecker) checkScope(scope string, h http.Header) error {
	return nil
}

type defaultScopeChecker struct{}

func (dsc *defaultScopeChecker) checkScope(scope string, h http.Header) error {
	token := extractToken(h.Get("Authorization"))

	ti, err := oauth2.RetrieveTokenInfo(tokenInfoService, token)
	if err != nil {
		msg := "Could not get token information. Authentication not possible."
		return errors.Wrap(err, msg)
	}

	return checkScopeList(scope, ti)
}

func extractToken(authorization string) (token string) {

	cutter := strings.Index(authorization, BEARER)
	if cutter >= 0 {
		token = authorization[cutter+6:] //6 == len("Bearer")
	}
	token = strings.Trim(token, " ")
	return
}

func checkScopeList(scope string, ti *oauth2.TokenInfo) error {
	hasScope := false
	for _, v := range ti.Scope {
		if v == scope {
			hasScope = true
			break
		}
	}
	if !hasScope {
		return errors.New("Scope not found! : " + scope)
	}
	return nil
}
