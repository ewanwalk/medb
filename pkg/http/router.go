package http

import (
	"encoder-backend/pkg/http/v1"
	"github.com/Ewan-Walker/respond"
	"github.com/gorilla/mux"
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

			if err, ok := data.(error); ok || status == http.StatusInternalServerError {

				if ok {
					data = err.Error()
				}

				return status, map[string]interface{}{
					"error": data,
					"code":  status,
				}
			}

			return status, map[string]interface{}{
				"code": status,
				"data": data,
			}
		},
	})

	return r
}

func (r *router) register() {

	// TODO serve static fs for assets
	// TODO obtain current "reports" of ongoing encodes

	v1.Register(r.Router)
}
