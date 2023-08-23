package rest

import (
	"encoding/json"
	"net/http"
)

//go:generate oapi-codegen -package rest -generate types,client,chi-server,spec -o ./handler.gen.go ./openapi.yaml

type PromotionHandler struct {
	PromotionsBackend PromotionsBackend
	AdminBackend      AdminBackend
}

var _ ServerInterface = (*PromotionHandler)(nil)

func NewPromotionHandler(pb PromotionsBackend, ad AdminBackend) PromotionHandler {
	return PromotionHandler{
		PromotionsBackend: pb,
		AdminBackend:      ad,
	}
}
func (h PromotionHandler) PostV1AdminProcessCsv(w http.ResponseWriter, r *http.Request) {
	var req ProcessCsv
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "failed", "Invalid request format")
		return
	}
	err := h.AdminBackend.Process(*req.FilePath)
	if err != nil {
		sendError(w, 500, "failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h PromotionHandler) PostV1AdminSwitchDb(w http.ResponseWriter, r *http.Request) {
	err := h.AdminBackend.SwitchStorage()
	if err != nil {
		sendError(w, 500, "failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h PromotionHandler) GetV1PromotionsId(w http.ResponseWriter, r *http.Request, id string) {

	out, err := h.PromotionsBackend.GetPromotion(id)

	if err != nil {
		sendError(w, 404, "not found", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Promotion{
		ExpirationDate: &out.ExpirationDate,
		Id:             &out.ID,
		Price:          &out.Price,
	})
}

func sendError(w http.ResponseWriter, code int, status, message string) {
	err := Error{
		Error:  &message,
		Status: &status,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
