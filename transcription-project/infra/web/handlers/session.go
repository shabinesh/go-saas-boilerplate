package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	AuthenticationKey = "authenticated"
)

func (h handlers) endSession(r *http.Request, w http.ResponseWriter) {
	session, _ := h.sessionStore.Get(r, "session")

	// Clear the session
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func (h handlers) getSession(r *http.Request, w http.ResponseWriter, vals map[string]interface{}) *sessions.Session {
	session, _ := h.sessionStore.Get(r, "session")

	for key, val := range vals {
		session.Values[key] = val
	}

	session.Save(r, w)

	return session
}
