//go:generate statik -f -src=../../frontend
package http

import (
	_ "encoder-backend/pkg/http/statik"
	"encoder-backend/pkg/http/utils"
	v1 "encoder-backend/pkg/http/v1"
	"github.com/ewanwalk/respond"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"net/http"
)

type router struct {
	*mux.Router
}

func newRouter() http.Handler {

	r := &router{
		Router: mux.NewRouter(),
	}

	r.register()

	respond.SetOptions(&respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {

			w.Header().Set("Access-Control-Allow-Origin", "*")

			if err, ok := data.(error); ok || status == http.StatusInternalServerError {

				if ok {
					data = err.Error()
				}

				return status, map[string]interface{}{
					"error": data,
					"code":  status,
				}
			}

			if dqr, ok := data.(utils.DQR); ok {
				return status, map[string]interface{}{
					"code":  status,
					"data":  dqr.Rows,
					"total": dqr.Total,
				}
			}

			return status, map[string]interface{}{
				"code": status,
				"data": data,
			}
		},
	})

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("Access-Control-Allow-Origin", "*")
		h.Add("Vary", "Origin")
		h.Add("Vary", "Access-Control-Request-Method")
		h.Add("Vary", "Access-Control-Request-Headers")
		h.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Request-With")
		h.Add("Access-Control-Allow-Methods", "GET,PUT,POST,OPTIONS")
	})

	return r
}

func (r *router) register() {

	// TODO serve static fs for assets
	// TODO obtain current "reports" of ongoing encodes
	v1.Register(r.Router)

	r.static()
}

func (r *router) static() {

	staticFS, err := fs.New()
	if err != nil {
		logrus.WithError(err).Fatal("http: failed to use vfs")
	}

	handler := http.FileServer(staticFS)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", handler))

	r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
}
