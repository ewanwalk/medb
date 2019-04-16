package v1

import (
	"encoder-backend/pkg/http/utils"
	"encoder-backend/pkg/models"
	"github.com/ewanwalk/gorm"
	"github.com/ewanwalk/respond"
	"github.com/gorilla/mux"
	"net/http"
)

func encodes(mux *mux.Router) {

	encodes := mux.PathPrefix("/encodes").Subrouter()
	encodes.HandleFunc("", getEncodes).Methods("GET")
	encodes.HandleFunc("/{encode}", getEncode).Methods("GET")

}

// getEncodes
// obtain all encodes based on the provided query filters
func getEncodes(w http.ResponseWriter, r *http.Request) {
	params := utils.Vars(r)

	encodes := make([]models.Encode, 0)

	if err := db.Scopes(models.Dyanmic(models.Encode{}, params)).Find(&encodes).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, encodes)
}

// getEncode
// obtain the provided encode
func getEncode(w http.ResponseWriter, r *http.Request) {
	params := utils.Vars(r)

	id, ok := params["encode"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	encode := models.Encode{}

	if err := db.Preload("File").Preload("QualityProfile").First(&encode, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, encode)
}
