package rest

import "promotions-service/internal"

type PromotionsBackend interface {
	GetPromotion(id string) (internal.PromotionRecord, error)
}

type AdminBackend interface {
	Process(file string) error
	SwitchStorage() error
}
