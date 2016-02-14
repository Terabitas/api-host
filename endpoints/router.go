package endpoints

import (
	"net/http"

	log "github.com/nildev/api-host/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/nildev/api-host/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/nildev/api-host/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/nildev/api-host/Godeps/_workspace/src/github.com/nildev/lib/router"
	"github.com/nildev/api-host/gen"
	"github.com/nildev/api-host/version"
)

var (
	ctxLog *log.Entry
)

func init() {
	ctxLog = log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)
}

func Router() http.Handler {
	r := mux.NewRouter()

	RegisterRoutes(r)

	return negroni.New(
		negroni.Wrap(r),
	)
}

func RegisterRoutes(r *mux.Router) {
	routes := gen.BuildRoutes()

	for _, rt := range routes {
		ctxLog.WithField("base-path", rt.BasePattern).Debugf("%+v", rt.Routes)
		r.PathPrefix(rt.BasePattern).Handler(negroni.New(
			negroni.Wrap(router.NewRouter(rt)),
		))
	}
}
