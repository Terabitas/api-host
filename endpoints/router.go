package endpoints

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nildev/api-host/config"
	"github.com/nildev/api-host/gen"
	"github.com/nildev/api-host/version"
)

var (
	ctxLog *log.Entry
)

func init() {
	ctxLog = log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)
}

func Router(cfg config.Config) http.Handler {
	r := mux.NewRouter()

	RegisterRoutes(r, cfg)

	return negroni.New(
		negroni.Wrap(r),
	)
}

func RegisterRoutes(r *mux.Router, cfg config.Config) {
	routes := gen.BuildRoutes()
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	for _, rt := range routes {
		apiRouter := mux.NewRouter().StrictSlash(true)
		for _, route := range rt.Routes {
			handlerFunc := negroni.New()
			// Add JWT middleware if route is protected
			if route.Protected {
				handlerFunc.Use(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext))
			}

			apiRouter.
				Methods(route.Method).
				Name(route.Name).
				Path(fmt.Sprintf("%s%s", rt.BasePattern, route.Pattern)).
				Queries(route.Queries...).
				HandlerFunc(route.HandlerFunc)

			// Add actual handler
			handlerFunc.Use(negroni.Wrap(apiRouter))
			r.PathPrefix(rt.BasePattern).Handler(handlerFunc)

			ctxLog.WithField("base-path", rt.BasePattern).
				WithField("protected", route.Protected).
				WithField("name", route.Name).
				WithField("path", fmt.Sprintf("%s%s", rt.BasePattern, route.Pattern)).
				WithField("method", route.Method).
				WithField("query", fmt.Sprintf("%+v", route.Queries)).
				Debugf("%+v", route)
		}
	}
}
