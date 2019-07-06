package v1

import (
	"encoder-backend/pkg/http/utils"
	"encoder-backend/pkg/models"
	"encoding/json"
	"github.com/ewanwalk/gorm"
	"github.com/ewanwalk/respond"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func profiles(mux *mux.Router) {

	profiles := mux.PathPrefix("/profiles").Subrouter()
	profiles.HandleFunc("", getProfiles).Methods("GET")
	profiles.HandleFunc("", createProfile).Methods("POST")
	profiles.HandleFunc("/{profile}", getProfile).Methods("GET")
	profiles.HandleFunc("/{profile}", updateProfile).Methods("PUT")
	profiles.HandleFunc("/{profile}/paths", getProfilePaths).Methods("GET")
	profiles.HandleFunc("/{profile}/encodes", getProfileEncodes).Methods("GET")

}

// getProfiles
// obtain all quality profiles based on filter parameters allowed
func getProfiles(w http.ResponseWriter, r *http.Request) {
	params := utils.Vars(r)

	profiles := make([]models.QualityProfile, 0)
	count := 0

	if err := db.Scopes(models.Dynamic(models.QualityProfile{}, params)).Find(&profiles).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Scopes(models.DynamicTotal(models.QualityProfile{}, params)).Count(&count).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, utils.DQR{Total: count, Rows: profiles})
}

// getProfile
// obtain a single quality profile based on the id provided
func getProfile(w http.ResponseWriter, r *http.Request) {
	params := utils.Vars(r)

	id, ok := params["profile"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	profile := models.QualityProfile{}

	if err := db.First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, profile)
}

// updateProfile
// modify a quality profile based on the request body provided
//
// Example:
// {
//    "name": "new name",
//    "threads": 0.25,
//    "codec: "x264
// }
func updateProfile(w http.ResponseWriter, r *http.Request) {

	params := utils.Vars(r)

	id, ok := params["profile"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	update := make(map[string]interface{})

	if err := json.Unmarshal(body, &update); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	profile := models.QualityProfile{}

	valid := models.Columns(profile)
	for key := range update {
		if _, ok := valid[key]; !ok {
			delete(update, key)
		}
	}

	// you cannot edit these fields
	for _, field := range []string{"encodes", "paths"} {
		delete(update, field)
	}

	if err := db.First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Model(&profile).Updates(update).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, profile)
}

// createProfile
// create a new quality profile based on the request body provided
//
// Example:
// {
//   "name": "my name"
// }
func createProfile(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	profile := models.QualityProfile{}

	if err := json.Unmarshal(body, &profile); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	except := profile.IsValid()
	if len(except) != 0 {
		respond.With(w, r, http.StatusBadRequest, except)
		return
	}

	if err := db.Create(&profile).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusCreated, profile)
}

// getProfilePaths
// obtain the paths which are currently utilizing the provided quality profile
func getProfilePaths(w http.ResponseWriter, r *http.Request) {

	params := utils.Vars(r)

	id, ok := params["profile"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	profile := models.QualityProfile{}
	count := 0

	if err := db.Preload("Paths", models.Dynamic(models.Path{}, utils.QueryParams(r))).
		First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Scopes(models.DynamicTotal(models.Path{}, params)).
		Where("quality_profile_id = ?", id).Count(&count).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, utils.DQR{Total: count, Rows: profile.Paths})
}

// getProfileEncodes
// obtain all the encodes which have used the quality profile provided
func getProfileEncodes(w http.ResponseWriter, r *http.Request) {

	params := utils.Vars(r)

	id, ok := params["profile"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	profile := models.QualityProfile{}
	count := 0

	if err := db.Preload("Encodes", models.Dynamic(models.Encode{}, utils.QueryParams(r))).
		First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Scopes(models.DynamicTotal(models.Encode{}, params)).
		Where("quality_profile_id = ?", id).Count(&count).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, utils.DQR{Total: count, Rows: profile.Encodes})
}
