package v1

import (
	"encoder-backend/pkg/database"
	"errors"
	"github.com/ewanwalk/gorm"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	db *gorm.DB
	//errDatabase       = errors.New("unexpected error has occurred")
	errInvalidRequest = errors.New("invalid request")
	errNotFound       = errors.New("resource not found")
)

func init() {
	dba, err := database.Connect()
	if err != nil {
		logrus.WithError(err).Fatalf("http.server.v1: failed to connect to database")
	}

	db = dba
}

// Register
// attach all routes which serve v1
func Register(mux *mux.Router) {

	v1 := mux.PathPrefix("/v1").Subrouter()

	// File related paths
	files(v1)
	// Encode related paths
	encodes(v1)
	// Paths
	paths(v1)
	// Profiles
	profiles(v1)

}
