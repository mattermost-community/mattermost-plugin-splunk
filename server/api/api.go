package api

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-plugin-splunk/server/config"
	"github.com/mattermost/mattermost-plugin-splunk/server/splunk"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/mlog"
)

const (
	// WebhookEndpoint a
	WebhookEndpoint = "/alert_action_wh"
)

// Error - returned error message for api errors
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// NewHTTPHandler initializes the router.
func NewHTTPHandler(sp splunk.Splunk, config *config.Config) http.Handler {
	return newHandler(sp, config)
}

// handler is an http.handler for all plugin HTTP endpoints
type handler struct {
	*mux.Router
	sp splunk.Splunk
}

type handlerWithUserID func(w http.ResponseWriter, r *http.Request, userID string)

func newHandler(sp splunk.Splunk, c *config.Config) *handler {
	h := &handler{
		Router: mux.NewRouter(),
		sp:     sp,
	}
	apiRouter := h.Router.PathPrefix(config.APIPath).Subrouter()

	apiRouter.HandleFunc(WebhookEndpoint, h.handleAlertActionWH(c)).Methods(http.MethodPost)

	return h
}

func (h *handler) handleAlertActionWH(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req splunk.AlertActionWHPayload
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.sp.LogError("Bad Request", "error", err.Error())
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			errMsg := "Bad webhook request. Missing url param 'id'"
			h.sp.LogError(errMsg)
			h.jsonError(w, Error{Message: errMsg, StatusCode: http.StatusBadRequest})
		}

		secret := r.URL.Query().Get("secret")
		if secret == "" {
			errMsg := "Bad webhook request. Missing url param 'secret'"
			h.sp.LogError(errMsg)
			h.jsonError(w, Error{Message: errMsg, StatusCode: http.StatusBadRequest})
			return
		}

		if secret != config.Secret {
			errMsg := "Bad webhook request. Invalid 'secret'"
			h.sp.LogError(errMsg)
			h.jsonError(w, Error{Message: errMsg, StatusCode: http.StatusBadRequest})
			return
		}

		err = h.sp.Notify(id, req)
		if err != nil {
			errMsg := "Error during webhook notify process"
			h.sp.LogError(errMsg, "error", err.Error())
			h.jsonError(w, Error{Message: errMsg, StatusCode: http.StatusInternalServerError})
			return
		}

		h.respondWithSuccess(w)
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
