package api

import (
	"errors"
	"net/http"
	"root/lib/errs"
	"root/lib/responser"
	"root/services/middleware"
)

func (h *Handlers) BuyItem(w http.ResponseWriter, r *http.Request) {

	u, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		responser.SendHttpError(w, "can't find user", http.StatusInternalServerError)
		return
	}

	item := r.PathValue("item")

	if item == "" {
		responser.SendHttpError(w, "wrong request, empty item", http.StatusBadRequest)
		return
	}

	err := h.shoppingService.Shop(u.ID, item)

	if err != nil && errors.Is(err, errs.NotEnoughMoney) {
		responser.SendHttpError(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		responser.SendHttpError(w, "problems with service, try later", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
