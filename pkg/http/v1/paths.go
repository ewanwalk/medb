package v1

import (
	"encoder-backend/pkg/http/utils"
	"encoder-backend/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/Ewan-Walker/gorm"
	"github.com/Ewan-Walker/respond"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func paths(mux *mux.Router) {

	paths := mux.PathPrefix("/paths").Subrouter()
	paths.HandleFunc("", getPaths).Methods("GET")
	paths.HandleFunc("", createPath).Methods("POST")
	paths.HandleFunc("/{path}", getPath).Methods("GET")
	paths.HandleFunc("/{path}", updatePath).Methods("PUT")
}

// getPaths
// obtain all paths based on the query filters provided
func getPaths(w http.ResponseWriter, r *http.Request) {

	params := utils.Vars(r)

	paths := make([]models.Path, 0)

	if err := db.Scopes(models.Dyanmic(models.Path{}, params)).Find(&paths).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, paths)
}

// getPath
// obtain the path of the provided id
func getPath(w http.ResponseWriter, r *http.Request) {
	params := utils.Vars(r)

	id, ok := params["path"]
	if !ok {
		respond.With(w, r, http.StatusBadRequest, errInvalidRequest)
		return
	}

	path := models.Path{}

	if err := db.Preload("QualityProfile").First(&path, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, path)
}

// updatePath
// edit fields of a path
//
// Example:
// {
//   "name": "my new name",
//   "event_scan_interval": 1000
// }
func updatePath(w http.ResponseWriter, r *http.Request) {

	params := utils.Vars(r)

	id, ok := params["path"]
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

	path := models.Path{}

	valid := models.Columns(path)
	for key := range update {
		if _, ok := valid[key]; !ok {
			delete(update, key)
		}
	}

	// you cannot edit these fields
	for _, field := range []string{"directory", "quality_profile"} {
		delete(update, field)
	}

	if err := db.First(&path, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.With(w, r, http.StatusNotFound, errNotFound)
			return
		}

		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Model(&path).Updates(update).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, path)
}

// createPath
// creates a new path based on the data provided
//
// Example:
// {
//    "name": "my new path",
//    "directory": "/mnt/movies/"
// }
func createPath(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	path := models.Path{}

	if err := json.Unmarshal(body, &path); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	except := path.IsValid()
	if len(except) != 0 {
		respond.With(w, r, http.StatusBadRequest, except)
		return
	}

	// we need to ensure no path duplicates exist & there are no nested paths
	count := 0
	db.Model(&path).Where("directory LIKE ?", fmt.Sprintf("%s%%", path.Directory)).Count(&count)
	if count != 0 {
		respond.With(w, r, http.StatusBadRequest, map[string]string{
			"directory": "The directory provided already exists or is a part of another path",
		})
		return
	}

	if err := db.Create(&path).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusCreated, path)
}
