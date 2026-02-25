package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"root/lib/responser"
	"root/services/middleware"

	"github.com/jackc/pgx/v5"
)

func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		responser.SendHttpError(w, "empty user", http.StatusInternalServerError)
		return
	}

	ua, err := h.userinfoService.Activity(u.ID)

	if err != nil && err == pgx.ErrNoRows {
		responser.SendHttpError(w, "user don't exists", http.StatusBadRequest)
		return
	} else if err != nil {
		responser.SendHttpError(w, "service error", http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(ua); err != nil {
		responser.SendHttpError(w, "service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
