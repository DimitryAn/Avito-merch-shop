package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"root/lib/errs"
	"root/lib/responser"
)

type requestData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {

	if head := r.Header.Get("Content-Type"); head != "application/json" {
		responser.SendHttpError(w, "wrong Content-Type", http.StatusBadRequest)
		return
	}

	var req requestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responser.SendHttpError(w, "empty body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	name, password := req.UserName, req.Password

	if name == "" || password == "" {
		responser.SendHttpError(w, "empty body", http.StatusBadRequest)
		return
	}

	jwt, err := h.loginService.Login(name, password)

	if errors.Is(err, errs.WrongPassword) {
		responser.SendHttpError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err != nil || jwt == "" {
		responser.SendHttpError(w, "can't make token", http.StatusInternalServerError)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{Token: jwt})
}
