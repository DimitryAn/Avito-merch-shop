package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"root/lib/errs"
	"root/lib/responser"
	"root/services/middleware"

	"github.com/jackc/pgx/v5"
)

func (h *Handlers) SendCoins(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		responser.SendHttpError(w, "serves error, empty user", http.StatusInternalServerError)
		return
	}

	if header := r.Header.Get("Content-Type"); header != "application/json" {
		responser.SendHttpError(w, "wrong content-type", http.StatusBadRequest)
		return
	}

	var reqData struct {
		Amount int    `json:"amount"`
		ToUser string `json:"toUser"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil || reqData.ToUser == "" || reqData.ToUser == u.Username || reqData.Amount < 0 {
		responser.SendHttpError(w, "wrong request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.sendCoinsService.SendCoins(u.ID, reqData.ToUser, reqData.Amount)

	if err != nil && errors.Is(err, errs.NotEnoughMoney) {
		responser.SendHttpError(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil && pgx.ErrNoRows == err {
		responser.SendHttpError(w, "user don't exists", http.StatusBadRequest)
		return
	} else if err != nil {
		responser.SendHttpError(w, "services error, try later", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
