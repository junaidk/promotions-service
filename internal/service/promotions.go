package service

import (
	"fmt"
	"promotions-service/internal"
)

type PermotionsRepository interface {
	Get(id string) *internal.PromotionRecord
}
type Promotions struct {
	repo PermotionsRepository
}

func NewPromotions(repo PermotionsRepository) Promotions {
	return Promotions{
		repo: repo,
	}
}
func (p *Promotions) GetPromotion(id string) (internal.PromotionRecord, error) {
	rec := p.repo.Get(id)

	if rec == nil {
		return internal.PromotionRecord{}, fmt.Errorf("id %s does not exist", id)
	}
	return *rec, nil
}
