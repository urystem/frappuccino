package service

import (
	"frappuccino/internal/dal"
	"frappuccino/models"
)

type aggregationService struct {
	aggreDalInter dal.AggregationDalInter
}

type AggregationServiceInter interface {
	SumOrder() (float64, error)
	PopularItems() ([]models.PopularItems, error)
}

func ReturnAggregationService(aggDalInter dal.AggregationDalInter) AggregationServiceInter {
	return &aggregationService{aggreDalInter: aggDalInter}
}

func (ser *aggregationService) SumOrder() (float64, error) {
	return ser.aggreDalInter.AmountSales()
}

func (ser *aggregationService) PopularItems() ([]models.PopularItems, error) {
	return ser.aggreDalInter.Popularies()
}
