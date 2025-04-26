package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	Search(find, from, minPrice, maxPrice string) (*models.SearchThings, error)
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

func (ser *aggregationService) Search(find, filter, minPrice, maxPrice string) (*models.SearchThings, error) {
	find = strings.Join(strings.Fields(find), " | ")
	if len(find) == 0 {
	}
	var minPriceF, maxPriceF float64
	var err error
	filtersMap := map[string]bool{"inventory": false, "menu": false, "orders": false}
	if len(minPrice) != 0 {
		minPriceF, err = strconv.ParseFloat(minPrice, 64)
		if err != nil {
		}
	}
	if len(maxPrice) != 0 {
		maxPriceF, err = strconv.ParseFloat(maxPrice, 64)
		if err != nil {
		} else if maxPriceF == 0 {
		}
	} else {
		maxPriceF = math.MaxFloat64
	}
	if len(filter) == 0 || filter == "all" {
		for k := range filtersMap {
			filtersMap[k] = true
		}
	} else {
		for _, from := range strings.FieldsFunc(filter, func(r rune) bool { return r == ',' || r == ' ' }) {
			if was, x := filtersMap[from]; !x {
				fmt.Println(from)
				return nil, errors.New("")
			} else if was {
				return nil, errors.New("")
			} else {
				filtersMap[from] = true
			}
		}
	}
	var ansSearch models.SearchThings
	for k, v := range filtersMap {
		if k == "inventory" && v {
			ser.aggreDalInter.SearchByWordInventory(find, minPriceF, maxPriceF, &ansSearch)
		}
	}
	return &ansSearch, nil
}

func (ser *aggregationService) timeParser(date string) (*time.Time, error) {
	if len(date) == 0 {
		return nil, nil
	}
	time, err := time.Parse("02.01.2006", date)
	return &time, err
}
