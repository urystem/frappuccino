package service

import (
	"fmt"
	"time"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type aggregationService struct {
	aggreDalInter dal.AggregationDalInter
}

type AggregationServiceInter interface {
	SumOrder() (float64, error)
	PopularItems() (*models.PopularItems, error)
	NumberOfOrderedItemsService(start, end string) (map[string]uint64, error)
}

func ReturnAggregationService(aggDalInter dal.AggregationDalInter) AggregationServiceInter {
	return &aggregationService{aggreDalInter: aggDalInter}
}

func (ser *aggregationService) SumOrder() (float64, error) {
	return ser.aggreDalInter.AmountSales()
}

func (ser *aggregationService) PopularItems() (*models.PopularItems, error) {
	return ser.aggreDalInter.Popularies()
}

func (ser *aggregationService) NumberOfOrderedItemsService(start, end string) (map[string]uint64, error) {
	startTime, err := ser.timeParser(start)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	endTime, err := ser.timeParser(end)
	if err != nil {
		return nil, err
	}

	return ser.aggreDalInter.CountOfOrderedItems(startTime, endTime)
}

func (ser *aggregationService) timeParser(date string) (*time.Time, error) {
	if len(date) == 0 {
		return nil, nil
	}
	time, err := time.Parse("02.01.2006", date)
	return &time, err
}
