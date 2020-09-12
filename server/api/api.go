package api

import (
	"encoding/json"
	"github.com/bakurits/mattermost-plugin-splunk/server/config"
	"net/http"

	"github.com/bakurits/mattermost-plugin-splunk/server/splunk"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/mlog"
)

// Error - returned error message for api errors
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// NewHTTPHandler initializes the router.
func NewHTTPHandler(sp splunk.Splunk) http.Handler {
	return newHandler(sp)
}

// handler is an http.handler for all plugin HTTP endpoints
type handler struct {
	*mux.Router
	sp splunk.Splunk
}

type handlerWithUserID func(w http.ResponseWriter, r *http.Request, userID string)

func newHandler(sp splunk.Splunk) *handler {
	h := &handler{
		Router: mux.NewRouter(),
		sp:     sp,
	}
	apiRouter := h.Router.PathPrefix(config.APIPath).Subrouter()

	apiRouter.HandleFunc("/alert_action_wh", h.handleAlertActionWH()).Methods(http.MethodPost)

	return h
}

func (h *handler) handleAlertActionWH() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req splunk.AlertActionWHPayload
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		h.sp.NotifyAll(req)
	}
}

func (h *handler) jsonError(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.StatusCode)
	h.respondWithJSON(w, err)
}

func (h *handler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		mlog.Error(err.Error())
	}
}

func (h *handler) respondWithSuccess(w http.ResponseWriter) {
	h.respondWithJSON(w, struct {
		Status string `json:"status"`
	}{Status: "OK"})
}

func (h *handler) extractUserIDMiddleware(handler handlerWithUserID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get("Mattermost-User-ID")
		if mattermostUserID == "" {
			h.jsonError(w, Error{Message: "Not Authorized", StatusCode: http.StatusUnauthorized})
			return
		}
		handler(w, r, mattermostUserID)
	}
}
