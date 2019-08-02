package okr2go

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/oxisto/go-httputil"
)

// NewRouter returns a configured mux router containing all REST endpoints
func NewRouter() *mux.Router {
	// pack angular ui
	box := packr.NewBox("./okr2go-ui/dist/okr2go-ui")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/objectives", GetObjectives)
	router.HandleFunc("/api/objectives/{objectiveID}/{resultID}/plus", ResultPlusOne)
	router.PathPrefix("/").Handler(http.FileServer(box))

	return router
}

func GetObjectives(w http.ResponseWriter, r *http.Request) {
	var err error

	objectives, err := ParseMarkdown("example.md")

	if err != nil {
		httputil.JSONResponse(w, r, nil, err)
		return
	}

	httputil.JSONResponse(w, r, objectives, err)
}

func ResultPlusOne(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		objectives []*Objective
		result     *KeyResult
	)

	// @todo Cache objectives in memory instead of loading them from markdown in every request
	objectives, err = ParseMarkdown("example.md")

	if err != nil {
		httputil.JSONResponse(w, r, nil, err)
		return
	}

	result, err = getResultFromRequest(objectives, w, r)
	if err != nil {
		httputil.JSONResponseWithStatus(w, r, nil, err, http.StatusBadRequest)
		return
	}

	if result == nil {
		httputil.JSONResponseWithStatus(w, r, nil, err, http.StatusNotFound)
		return
	}

	if result.Current == result.Target {
		httputil.JSONResponseWithStatus(w, r, nil, nil, http.StatusNotModified)
		return
	}

	result.Current++

	// @todo Persist changes

	httputil.JSONResponseWithStatus(w, r, result, nil, http.StatusOK)
	return
}

func getObjectiveFromRequest(objectives []*Objective, w http.ResponseWriter, r *http.Request) (objective *Objective, err error) {
	var (
		ok                bool
		objectiveID       int
		objectiveIDString string
	)

	if objectiveIDString, ok = mux.Vars(r)["objectiveID"]; !ok {
		return nil, errors.New("Request did not contain a resultID")
	}

	if objectiveID, err = strconv.Atoi(objectiveIDString); err != nil {
		return nil, errors.New("Could not parse objectiveID")
	}

	if objectiveID < 0 || objectiveID > len(objectives) {
		return nil, nil
	}

	objective = objectives[objectiveID]

	return objective, nil
}

func getResultFromRequest(objectives []*Objective, w http.ResponseWriter, r *http.Request) (result *KeyResult, err error) {
	var (
		ok        bool
		objective *Objective
		resultID  string
	)

	// return if we either have an error (which we pass down, resulting in a BadRequest)
	// or if the objective was not found (resulting in a NotFound)
	if objective, err = getObjectiveFromRequest(objectives, w, r); err != nil || objective == nil {
		return nil, err
	}

	if resultID, ok = mux.Vars(r)["resultID"]; !ok {
		return nil, errors.New("Request did not contain a resultID")
	}

	result = objective.FindKeyResult(resultID)

	return result, nil
}
