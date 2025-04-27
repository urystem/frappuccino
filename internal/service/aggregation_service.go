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
	OrderedItemsPeriod(period, month, year string) (*models.OrderStats, error)
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

	filtersMap := map[string]bool{"inventory": false, "menu": false, "orders": false}
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
		if v {
			var count uint64
			switch k {
			case "inventory":
				err = ser.aggreDalInter.SearchByWordInventory(find, minPriceF, maxPriceF, &ansSearch)
				ansSearch.Inventory_math = new(uint64)
				count = uint64(len(ansSearch.Inventories))
				*ansSearch.Inventory_math = count
			case "menu":
				err = ser.aggreDalInter.SearchByWordMenu(find, minPriceF, maxPriceF, &ansSearch)
				ansSearch.Menu_math = new(uint64)
				count = uint64(len(ansSearch.Menus))
				*ansSearch.Menu_math = count
			case "orders":
				err = ser.aggreDalInter.SearchByWordOrder(find, minPriceF, maxPriceF, &ansSearch)
				ansSearch.Order_math = new(uint64)
				count = uint64(len(ansSearch.Orders))
				*ansSearch.Order_math = count
			}
			if err != nil {
				return nil, err
			} else {
				ansSearch.Total_math += count
			}
		}
	}
	return &ansSearch, nil
}

func (ser *aggregationService) OrderedItemsPeriod(period, month, year string) (*models.OrderStats, error) {
	var orderStats models.OrderStats
	var err error
	switch strings.ToLower(period) {
	case "day":
		mouthInt := time.Now().Month()

		if len(month) != 0 {
			monthTime, err := time.Parse("January", month)
			if err != nil {
				return nil, err
			}
			// mouthInt = monthTime.UTC().Month()
			mouthInt = monthTime.Month()
		}
		orderStats.OrderItems, err = ser.aggreDalInter.PeriodMonth(mouthInt)
		if err != nil {
			return nil, err
		}
		orderStats.Month = mouthInt.String()
	case "month":
		yearInt := time.Now().Year()
		if len(year) != 0 {
			yearAtoi, err := strconv.Atoi(year)
			if err != nil {
				return nil, err
			} else if yearAtoi < 2000 || yearInt < yearAtoi {
				return nil, errors.New("dd")
			}
			yearInt = yearAtoi
		}
		orderStats.OrderItems, err = ser.aggreDalInter.PeriodYear(yearInt)
		if err != nil {
			return nil, err
		}
		orderStats.Year = yearInt
	default:
		return nil, errors.New("dsf")
	}
	orderStats.Period = period
	return &orderStats, nil
}

func (ser *aggregationService) timeParser(date string) (*time.Time, error) {
	if len(date) == 0 {
		return nil, nil
	}
	time, err := time.Parse("02.01.2006", date)
	return &time, err
}
